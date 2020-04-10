package main

import (
	"context"
	"github.com/k2rth1k/grpc_learning/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("[Greet function is called with request %v]\n",req)
	firstName := req.Greeting.FirstName
	result := "Hello " + firstName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func(*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest,stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Println("[called GreetManyTimes rpc....]")
	firstName:= req.Greeting.FirstName
	for i:=0;i<10;i++{
		result:="Hello " + firstName +" number-" +strconv.Itoa(i)
		res:=&greetpb.GreetManyTimesResponse{
			Result:               result,
		}
		_=stream.Send(res)
		time.Sleep(1000*time.Millisecond)
	}
	return nil
}

func main() {
	log.Println("[server is online .....]")

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
