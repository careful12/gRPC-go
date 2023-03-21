package main

import (
	pb "grpc-go/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	names := &pb.NamesList{
		Names: []string{"Alice", "Bob", "Joe"},
	}

	// Unary
	// callSayHello(client)

	// Server stream
	// callSayHelloServerStream(client, names)

	// Client stream
	// callSayHelloClientStream(client, names)

	// Bi-directional stream
	callSayHelloBidirectionalStream(client, names)
}
