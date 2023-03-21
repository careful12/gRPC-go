# gRPC-go

ref: 
- [AkhilSharma90 - FULL PROJECT - GO + GRPC](https://www.youtube.com/watch?v=a6G5-LUlFO4&ab_channel=AkhilSharma)
- [gRPC - Quick Start](https://grpc.io/docs/languages/go/quickstart/#update-and-run-the-application)
- [gRPC - Basics tutorial](https://grpc.io/docs/languages/go/basics/)

## Simply introduce types of gRPC (4 type)

### Unary
- Similar to traditional API, like RESTful
- **Client send request** and **Server send response**

### Server streaming
- Client send request
- Server send stream(data flow)
### Client streaming
- Client send stream(data flow)
- Server send request

### Bi-directional streaming
- Both client and server send stream

## Building Environment
* My notebook is windows 10
* libprotoc 3.19.4
    * [Protocol Buffer Compiler Installation - Install pre-compiled binaries (any OS)](https://grpc.io/docs/protoc-installation/#install-pre-compiled-binaries-any-os)
    * [Notes on installing GRPC development environment for Go on Windows](https://gist.github.com/jjeffery/6e3a4d18ffbe1fc715be403f6391c5f4)
    * You can test by run `protoc --version`
    * Success
![](https://i.imgur.com/8ewsjSM.png)

* Install the protocol compiler plugins for Go using the following commands:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
* run `go mod tidy`



## Creating proto files
```protobuf
syntax="proto3";

// Files will get created in thin path
option go_package = "./proto";

package greet_service;

service GreetService{
    // Unary
    rpc SayHello(NoParam) returns (HelloResponse);
    
    // Server stream 
    rpc SayHelloServerStreaming(NamesList) returns (stream HelloResponse);
    
    // Client stream
    rpc SayHelloClientStreaming(stream HelloRequest) returns (MessagesLists);

    // Bi-directional stream
    rpc SayHelloBidirectionalStreaming(stream HelloRequest) returns (stream HelloResponse);

}

message NoParam{};

message HelloRequest{
    string name = 1;
}

message HelloResponse{
    string message = 1 ;
}

message NamesList{
    repeated string names = 1;
}

message MessagesLists{
    repeated string messages = 1;
}
```
- proto files will be created when you run
`protoc --go_out=. --go-grpc_out=. proto/greet.proto`

## main.go for Server and Client

### Server
```go
// server/main.go
package main

import (
	"log"
	"net"
	pb "grpc-go/proto"
	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type helloServer struct {
	pb.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})
	log.Printf("Server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}

```

### Client 
```go
// client/main.go
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

        // If testing Unary , comment out this
	names := &pb.NamesList{
		Names: []string{"Alice", "Bob", "Joe"},
	}

        // Discomment out to test    
	// Unary
	// callSayHello(client)

	// Server stream
	// callSayHelloServerStream(client, names)

	// Client stream
	// callSayHelloClientStream(client, names)

	// Bi-directional stream
	// callSayHelloBidirectionalStream(client, names)
}

```

## Unary

![](https://i.imgur.com/FZfjTDl.png)

### Server
```go
// server/unary.go
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
```
### Client
```go
// client/unary.go
package main

import (
	"context"
	pb "grpc-go/proto"
	"log"
	"time"
)

func callSayHello(client pb.GreetServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SayHello(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("%s", res.Message)
}
```

## Server Stream
![](https://i.imgur.com/cLx4z9d.png)

### Server
```go
// server/server_stream.go
package main

import (
	"log"
	"time"
	pb "grpc-go/proto"
)

func (s *helloServer) SayHelloServerStreaming(req *pb.NamesList, stream pb.GreetService_SayHelloServerStreamingServer) error {
	log.Printf("Got request with names: %v", req.Names)
        // Send response stream    
	for _, name := range req.Names {
		res := &pb.HelloResponse{
			Message: "Hello " + name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}
```
### Client
```go
// client/server_stream.go
package main

import (
	"context"
	"io"
	"log"
	pb "grpc-go/proto"
)

func callSayHelloServerStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Streaming Started")
        // Send request    
	stream, err := client.SayHelloServerStreaming(context.Background(), names)
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}
	for {
               // Receive response stream           
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming %v", err)
		}
		log.Println(message)
	}
	log.Printf("Streaming finished.")
}

```

## Client Stream
![](https://i.imgur.com/yWcnS8O.png)

### Server
```go
// server/client_stream.go
package main

import (
	pb "grpc-go/proto"
	"io"
	"log"
)

func (s *helloServer) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {
	var messages []string
	for {
               // Receive request stream     
		req, err := stream.Recv()
		if err == io.EOF {
                       // Send response        
			return stream.SendAndClose(&pb.MessagesLists{Messages: messages})
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name: %v", req.Name)
		messages = append(messages, "Hello", req.Name)
	}
}
```
### Client
```go
// client/client_stream.go
package main

import (
	"context"
	pb "grpc-go/proto"
	"log"
	"time"
)

func callSayHelloClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Client streaming started.")
	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}
        // Send request stream
	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		log.Printf("Sent the request with name: %s", name)
		time.Sleep(2 * time.Second)
	}
        // Receive response
	res, err := stream.CloseAndRecv()
	log.Printf("Client streaming finished")
	if err != nil {
		log.Fatalf("Error while receiving %v", err)
	}
	log.Printf("%v", res.Messages)
}
```

## Bi-directional stream
![](https://i.imgur.com/Lyj7cfO.png)

### Server
```go
// server/bi_stream.go
package main

import (
	pb "grpc-go/proto"
	"io"
	"log"
)

func (s *helloServer) SayHelloBidirectionalStreaming(stream pb.GreetService_SayHelloBidirectionalStreamingServer) error {
	for {
		// Receive request stream
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name: %v", req.Name)

		// Send response stream
		res := &pb.HelloResponse{
			Message: "Hello " + req.Name,
		}

		if err := stream.Send(res); err != nil {
			return err
		}
	}
}
```
### Client
```go
// client/bi_stream.go
package main

import (
	"context"
	pb "grpc-go/proto"
	"io"
	"log"
	"time"
)

func callSayHelloBidirectionalStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Bidirectional streaming started.")
	stream, err := client.SayHelloBidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	// Receive response stream
	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming %v", err)
			}
			log.Println(message)

		}
	}()

	// Send request stream
	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		time.Sleep(2 * time.Second)
	}
	stream.CloseSend()

	log.Printf("Bidirectional streaming finished.")
}
```
