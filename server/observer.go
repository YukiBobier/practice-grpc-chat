package main

import (
	"sync"
)

func newPublisher[T any]() *publisher[T] {
	return &publisher[T]{subscribers: make(map[*subscriber[T]]chan *T)}
}

type publisher[T any] struct {
	subscribers map[*subscriber[T]]chan *T
	sync.Mutex
}

func (p *publisher[T]) addSubscriber(s *subscriber[T]) {
	p.Lock()
	p.subscribers[s] = s.ch
	p.Unlock()
}

func (p *publisher[T]) removeSubscriber(s *subscriber[T]) {
	p.Lock()
	delete(p.subscribers, s)
	p.Unlock()
}

func (p *publisher[T]) do(item *T) {
	for _, ch := range p.subscribers {
		ch := ch
		go func() {
			ch <- item
		}()
	}
}

func newSubscriber[T any]() *subscriber[T] {
	return &subscriber[T]{ch: make(chan *T)}
}

type subscriber[T any] struct {
	ch        chan *T
	publisher *publisher[T]
}

func (s *subscriber[T]) do(p *publisher[T]) {
	p.addSubscriber(s)
	s.publisher = p
}

func (s *subscriber[T]) close() {
	s.publisher.removeSubscriber(s)
	s.publisher = nil
	close(s.ch)
}
