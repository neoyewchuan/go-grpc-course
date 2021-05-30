package main

import (
	"context"
	"fmt"
	"log"

	"github.com/neoyewchuan/go-grpc-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnarySum(c)
	doUnarySub(c)
	doUnaryMul(c)
	doUnaryDiv(c)
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
