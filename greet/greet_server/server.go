package main

import (
	"context"
	"fmt"
	"github.com/k2rth1k/grpc_learning/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("[Greet function is called with request %v]\n", req)
	firstName := req.Greeting.FirstName
	result := "Hello " + firstName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Println("[called GreetManyTimes rpc....]")
	firstName := req.Greeting.FirstName
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number-" + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		_ = stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Println("[called LongGreet rpc]")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{Result: result})
		}
		if err != nil {
			log.Fatalf("[failed to recieve req:%v...]", err)
		}
		result = result + req.Greeting.FirstName + "!\n"
	}
	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Printf("[GreetEveryone RPC is called.....]")
	for {
		req, err := stream.Recv()
		if err != nil {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream:%v", err)
			return err
		}
		firstName := req.Greeting.FirstName
		result := "Hello " + firstName + "! "
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	log.Println("[Greet with Deadline is called]")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("the client canceled the request!")
			return nil, status.Error(codes.Canceled, "canceled by client")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.Greeting.FirstName
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{Result: result}
	return res, nil
}

func main() {
	log.Println("[server is online .....]")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	tls := true
	var s *grpc.Server
	if !tls || os.Args[0] == "--insecure" {
		s = grpc.NewServer()
	} else {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("[failed loading certificates:%v]", sslErr)
			return
		}
		opts := grpc.Creds(creds)
		s = grpc.NewServer(opts)
	}

	greetpb.RegisterGreetServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
