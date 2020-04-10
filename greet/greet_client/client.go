package main

import (
	"context"
	"fmt"
	"github.com/k2rth1k/grpc_learning/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//doUnary(c)
	//doStreaming(c)
	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient){
	fmt.Println("Starting to do Unary RPC.....")
	req := &greetpb.GreetRequest{
		Greeting:            &greetpb.Greeting{
			FirstName:            "karthik",
			LastName:             "chowdary",
		},
	}
	res,err:=c.Greet(context.Background(),req)
	if err!=nil{
		log.Fatalf("error while calling Greet RPC:%v",err)
	}
	log.Printf("Response from Greet: %v",res.Result)
	fmt.Printf("Created a client: %v", c)
}

func doStreaming(c greetpb.GreetServiceClient){
	res,err:=c.GreetManyTimes(context.Background(),&greetpb.GreetManyTimesRequest{Greeting:&greetpb.Greeting{
		FirstName:            "karthik",
		LastName:             "chowdary",

	}})
	if err!=nil{
		log.Fatalf("[error while calling GretManyTimes RPC:%v...]",err)
	}
	for {
		msg,err:=res.Recv()
		if err == io.EOF{
			break
		}
		if err!=nil{
			log.Fatalf("[error occurred while streaming: %v....]",err)
		}
		log.Printf("[Response from GreetManyTimes: %v....]",msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient){
	stream,err:=c.LongGreet(context.Background())
	if err!=nil{
		log.Fatalf("[error while calling LongGreet:%v]",err)
	}
	for i:=0;i<5;i++{
		req:=&greetpb.LongGreetRequest{Greeting:&greetpb.Greeting{FirstName:"karthik"}}
		err=stream.Send(req)
		if err!=nil{
			log.Fatalf("[failed to send request due to error:%v]",err)
		}
		log.Printf("[sent request:%v]",req)
		time.Sleep(1000*time.Millisecond)
	}
	res,err:=stream.CloseAndRecv()
	if err!=nil{
		log.Fatalf("[error while recieveing message from LongGreet:%v]",err)
	}
	log.Printf("Response fromLong Greet:%v",res)
}