package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/neoyewchuan/go-grpc-course/primenumberdecomposition/primenumdecompb"
)

func main() {
	fmt.Println("Hello I'm a client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := primenumdecompb.NewPrimeDecomServiceClient(cc)
	//fmt.Printf("Created client: %f", c)
	//doUnary(c)

	doServerStreaming(c)
}

func doServerStreaming(c primenumdecompb.PrimeDecomServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &primenumdecompb.PrimeDecomManyTimesRequest{
		Primedecom: &primenumdecompb.PrimeDecom{
			Number: 510510,
		},
	}
	resStream, err := c.PrimeDecomManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeDecomManyTimes: %v", msg.GetResult())
	}
}
