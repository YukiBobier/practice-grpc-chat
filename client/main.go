package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/YukiBobier/practice-grpc-chat/chat"
	"github.com/marcusolsson/tui-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	name    = flag.String("name", "no-name", "Your name")
	address = flag.String("address", "localhost:50051", "The address of gRPC server")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	client := pb.NewChatServiceClient(conn)

	runUI(client, *name)
}

func runUI(client pb.ChatServiceClient, name string) {
	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	root := tui.NewVBox(historyBox, inputBox)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatalf("Failed to create ui: %v", err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	input.OnSubmit(func(e *tui.Entry) {
		post(client, &pb.Message{Name: name, Body: e.Text()})
		input.SetText("")
	})

	// Keep receiving messages and listing them
	stream := subscribe(client, &pb.SubscribeRequest{})
	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to reading stream: %v", err)
			}
			ui.Update(func() {
				history.Append(tui.NewLabel(fmt.Sprintf(
					"%s %-10s %s",
					message.GetPostedAt().AsTime().Format(time.RFC3339),
					message.GetName(),
					message.GetBody(),
				)))
			})
		}
	}()

	if err := ui.Run(); err != nil {
		log.Fatalf("Failed to run ui: %v", err)
	}
}

func post(client pb.ChatServiceClient, message *pb.Message) {
	_, err := client.Post(context.Background(), message)
	if err != nil {
		log.Fatalf("Failed to call Post: %v", err)
	}
}

func subscribe(client pb.ChatServiceClient, request *pb.SubscribeRequest) pb.ChatService_SubscribeClient {
	stream, err := client.Subscribe(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to call Subscribe: %v", err)
	}
	return stream
}
