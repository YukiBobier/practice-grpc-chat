package main

import (
	"context"
	"log"

	pb "github.com/YukiBobier/practice-grpc-chat/internal/chat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newChatServiceServer() *chatServiceServer {
	return &chatServiceServer{publisher: newPublisher[pb.Message]()}
}

type chatServiceServer struct {
	publisher *publisher[pb.Message]
	pb.UnimplementedChatServiceServer
}

func (s *chatServiceServer) Post(ctx context.Context, message *pb.Message) (*pb.PostResponse, error) {
	log.Printf("`Post` is called: message (%v)\n", message)

	if err := ctx.Err(); err != nil {
		log.Printf("`Post` is cancelled: %v", err)
		return nil, nil
	}

	message.PostedAt = timestamppb.Now()
	s.publisher.do(message)

	return &pb.PostResponse{}, nil
}

func (s *chatServiceServer) Subscribe(req *pb.SubscribeRequest, stream pb.ChatService_SubscribeServer) error {
	log.Print("`Subscribe` is called.\n")

	ctx := stream.Context()
	if err := ctx.Err(); err != nil {
		log.Printf("`Subscribe` is cancelled: %v", err)
		return nil
	}

	subscriber := newSubscriber[pb.Message]()
	subscriber.do(s.publisher)
	defer subscriber.close()

	go func() {
		log.Print("The subscription is started.\n")
		for message := range subscriber.ch {
			if err := stream.Send(message); err != nil {
				log.Printf("Failed to send message (%v): %v", message, err)
			}
		}
		log.Print("The subscription is finished.\n")
	}()

	<-ctx.Done()
	log.Print("`Subscribe` is cancelled.")
	return nil
}
