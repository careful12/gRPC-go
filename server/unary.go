package main

import (
	"context"
	pb "grpc-go/proto"
)

func (s *helloServer) SayHello(ctx context.Context, req *pb.NoParam) (*pb.HelloResponse, error) {
	// Send response
	return &pb.HelloResponse{
		Message: "Hello",
	}, nil
}
