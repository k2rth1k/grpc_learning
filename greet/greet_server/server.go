package main

import (
	"context"
	"fmt"
	"github.com/k2rth1k/grpc_learning/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("[Greet function is called with request %v]\n",req)
	firstName := req.Greeting.FirstName
	result := "Hello " + firstName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
