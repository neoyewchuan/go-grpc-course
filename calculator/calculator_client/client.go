package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/neoyewchuan/go-grpc-course/calculator/calculatorpb"
)

func main() {
	fmt.Println("Hello I'm a client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnarySum(c)
	// doUnarySub(c)
	// doUnaryMul(c)
	// doUnaryDiv(c)

	//doServerStreaming(c)
	doClientStreaming(c)
}

func doUnarySum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC (Sum) ...")
	req := &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			FirstNumber:  242,
			SecondNumber: 235,
		},
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)
}

func doUnarySub(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC (Sub) ...")
	req := &calculatorpb.SubRequest{
		Sub: &calculatorpb.Sub{
			FirstNumber:  242,
			SecondNumber: 235,
		},
	}
	res, err := c.Sub(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sub RPC: %v", err)
	}
	log.Printf("Response from Sub: %v", res.SubResult)
}

func doUnaryMul(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC (Div) ...")
	req := &calculatorpb.MulRequest{
		Mul: &calculatorpb.Mul{
			FirstNumber:  242,
			SecondNumber: 235,
		},
	}
	res, err := c.Mul(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Mul RPC: %v", err)
	}
	log.Printf("Response from Mul: %v", res.MulResult)
}

func doUnaryDiv(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC (Div) ...")
	req := &calculatorpb.DivRequest{
		Div: &calculatorpb.Div{
			FirstNumber:  242,
			SecondNumber: 235,
		},
	}
	res, err := c.Div(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Div RPC: %v", err)
	}
	log.Printf("Response from Div: %v", res.DivResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeNumberDecomposition Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 510510,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	for {
		recv, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("%v", recv.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		{Number: 1},
		{Number: 11},
		{Number: 21},
		{Number: 33},
		{Number: 41},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling Computeaverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	log.Printf("ComputeAverage Response: %v\n", res.GetResult())
}
