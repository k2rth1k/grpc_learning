package main

import (
	"context"
	"github.com/k2rth1k/grpc_learning/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)
func(s *server) Sum(ctx context.Context,req *calculatorpb.SumRequest) (*calculatorpb.SumResponse,error){
	log.Printf("Sum api is called with req:%v.....",req)
	sum:=req.FirstNumber + req.SecondNumber
	return &calculatorpb.SumResponse{SumResult:sum},nil
}

func(s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error{
	log.Println("ComputeAverage rpc is called...")
	average:=0.0
	count:=-1
	for{
		count+=1
		req,err:=stream.Recv()
		if err==io.EOF{
			 _=stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				 Average:              average/float64(count),
			 })
			 return nil
		}
		if err!=nil{
			log.Fatalf("[failed to recaieve stream of request due to error:%v]",err)
		}
		average+=float64(req.Number)
	}
	return nil
}

type server struct{}
func main(){
	lis,err:=net.Listen("tcp","0.0.0.0:50051")
	if err!=nil{
		log.Fatalf("failed to start listener:%v\n",err)
	}
	s:=grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s,&server{})

	log.Println("[server is online.....]")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
