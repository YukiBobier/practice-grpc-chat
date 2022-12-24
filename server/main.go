package main

import (
	"flag"
	"log"
	"net"

	pb "github.com/YukiBobier/practice-grpc-chat/chat"
	"google.golang.org/grpc"
)

var address = flag.String("address", "localhost:50051", "The server address")

func main() {
	log.Print("The gRPC server is starting...")

	flag.Parse()

	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, newChatServiceServer())

	log.Printf("The gRPC server is serving on %s.\n", *address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
