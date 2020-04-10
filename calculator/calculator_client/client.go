package main

import (
	"context"
	"github.com/k2rth1k/grpc_learning/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main(){
	cc,err:=grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err!=nil{
		log.Fatalln("failed to dial or make connection")
	}

	c:=calculatorpb.NewCalculatorServiceClient(cc)
	req:=&calculatorpb.SumRequest{FirstNumber:1,SecondNumber:2}
	log.Printf("[ sending request to get sum with request:%v..... ]",req)
	response,err:=c.Sum(context.Background(),req)


	if err!=nil{
		log.Fatalf("failed to get Sum")
	}
	log.Printf("[response from Sum:%v...\n]",response)

	doClientSideStreaming(c)
}

func doClientSideStreaming(c calculatorpb.CalculatorServiceClient){
	log.Println("streaming numbers for  average")
	res,err:=c.ComputeAverage(context.Background())
	if err!=nil{
		log.Fatalf("[failed to call ComputeAverage RPC due to error:%v]",err)
	}
	for i:=0;i<10;i++{
		req:=&calculatorpb.ComputeAverageRequest{Number:int32(i+1)}
		err=res.Send(req)
		if err!=nil{
			log.Fatalf("[failed to send request due to error:%v]",err)
		}
		log.Printf("[sent request:%v]",req)
	}
	response,err:=res.CloseAndRecv()
	if err!=nil{
		log.Fatalf("[failed to recieve response from ComputeAverage]")
	}
	log.Printf("[Computeed Average of all numbers is:%v]",response.Average)
}