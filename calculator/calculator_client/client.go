package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doErrorUnarySqrt(c)
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

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v\n", err)
		return
	}

	waitc := make(chan struct{})

	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		numbers := []int64{1, 5, 13, 6, 2, 20, 12, 9, 33, 23}
		for _, number := range numbers {
			fmt.Printf("Sending number: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{Number: number})
			time.Sleep(500 * time.Millisecond)

		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages fromt the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving stream: %v", err)
				break
			}
			fmt.Printf("A New Maximum Discovered: %v\n", res.GetMaximumNumber())
			time.Sleep(500 * time.Millisecond)
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
}

func doErrorUnarySqrt(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Error Handling Unary RPC (SquareRoot) ...")

	// correct call
	calSquareRoot(c, 36)
	// correct call
	calSquareRoot(c, 235)
	// error call
	calSquareRoot(c, -23)
}

func calSquareRoot(c calculatorpb.CalculatorServiceClient, number int64) {

	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: number})
	if err != nil {
		respErr, userErr := status.FromError(err)
		if userErr {
			// actual error from gRPC (user error)
			fmt.Printf("Error message from server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("You probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
		}
	} else {
		fmt.Printf("SquareRoot of number %v is %v\n", number, res.GetNumberRoot())
	}
}
