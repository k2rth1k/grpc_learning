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
	log.Printf("[response from Sum:%v...]",response)
}
