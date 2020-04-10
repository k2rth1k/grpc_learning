package main

import (
	"context"
	"github.com/k2rth1k/grpc_learning/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)
func(s *server) Sum(ctx context.Context,req *calculatorpb.SumRequest) (*calculatorpb.SumResponse,error){
	log.Printf("Sum api is called with req:%v.....",req)
	sum:=req.FirstNumber + req.SecondNumber
	return &calculatorpb.SumResponse{SumResult:sum},nil
}

type server struct{}
func main(){
	lis,err:=net.Listen("tcp","0.0.0.0:50051")
	if err!=nil{
		log.Fatalf("failed to start listener:%v\n",err)
	}
	s:=grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s,&server{})

	log.Println("server is online.....")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
