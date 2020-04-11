package main

import (
	"context"
	"fmt"
	"github.com/k2rth1k/grpc_learning/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	log.Println("[Hello I'm a client]")

	tls:=true
	var opts grpc.DialOption
	if !tls || os.Args[0]=="--insecure"{
		opts=grpc.WithInsecure()
	}else{
		creds,sslErr:=credentials.NewClientTLSFromFile("ssl/ca.crt","")
		if sslErr!=nil{
			log.Fatalf("[failed loading certificates:%v]",sslErr)
			return
		}
		opts=grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//doUnary(c)
	//doStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	DoUnaryWithDeadline(c)
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

func doBiDiStreaming(c greetpb.GreetServiceClient){
	fmt.Println("Starting to do a BIDI Streaming RPC...")

	stream,err:=c.GreetEveryone(context.Background())
	if err!=nil{
		log.Fatalf("Errot while creating stream: %v",err)
		return
	}
	requests:=[]*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{Greeting:&greetpb.Greeting{FirstName:"karthik",LastName:"chowdary"}},
		&greetpb.GreetEveryoneRequest{Greeting:&greetpb.Greeting{FirstName:"john",LastName:"cena"}},
		&greetpb.GreetEveryoneRequest{Greeting:&greetpb.Greeting{FirstName:"john",LastName:"wick"}},
		&greetpb.GreetEveryoneRequest{Greeting:&greetpb.Greeting{FirstName:"walter",LastName:"white"}},
		&greetpb.GreetEveryoneRequest{Greeting:&greetpb.Greeting{FirstName:"arthur",LastName:"morgan"}},
	}

	waitc:=make(chan struct{})
	go func(){
		for _, req := range requests{
			err=stream.Send(req)
			if err!=nil{
				log.Fatalf("[failed to send message due to error:%v]",err)
			}
			log.Printf("[sucessfully sent req:%v]\n",req)
			time.Sleep(100*time.Millisecond)
		}
		err=stream.CloseSend()
		if err!=nil{
			log.Fatalf("[failed to close straming:%v]\n",err)
		}
	}()

	go func(){
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while recieving: %v", err)
				break
			}
			log.Printf("[Recieved: %v]\n", res.Result)
		}
		close(waitc)
	}()
	<-waitc
}

func DoUnaryWithDeadline(c greetpb.GreetServiceClient){
	fmt.Println("[starting to do a UnaryWithDeadline RPC...]")
	req:=&greetpb.GreetWithDeadlineRequest{Greeting:&greetpb.Greeting{FirstName:"karthik",LastName:"chowdary"}}
	ctx,cancelFunc:=context.WithTimeout(context.Background(),4*time.Second)
	res,err:=c.GreetWithDeadline(ctx,req)
	defer cancelFunc()
	if err!=nil{
		log.Fatalf("error while calling GreetWithDeadline RPC: %v",err)
	}
	log.Printf("Response from GreetWithDeadline: %v",res.Result)
}