package main

import (
	"context"
	"fmt"
	"log"

	"github.com/neoyewchuan/go-grpc-course/simplesum/sumpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := sumpb.NewSumServiceClient(cc)

	doUnary(c)
}

func doUnary(c sumpb.SumServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	req := &sumpb.SumRequest{
		Sum: &sumpb.Sum{
			Value_1: 242,
			Value_2: 235,
		},
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Result)
}
