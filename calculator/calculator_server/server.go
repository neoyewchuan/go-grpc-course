package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/neoyewchuan/go-grpc-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with: %v\n", req)
	value1 := req.GetSum().GetFirstNumber()
	value2 := req.GetSum().GetSecondNumber()
	result := value1 + value2
	res := &calculatorpb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func (*server) Sub(ctx context.Context, req *calculatorpb.SubRequest) (*calculatorpb.SubResponse, error) {
	fmt.Printf("Sub function was invoked with: %v\n", req)
	value1 := req.GetSub().GetFirstNumber()
	value2 := req.GetSub().GetSecondNumber()
	result := value1 - value2
	res := &calculatorpb.SubResponse{
		SubResult: result,
	}
	return res, nil
}

func (*server) Mul(ctx context.Context, req *calculatorpb.MulRequest) (*calculatorpb.MulResponse, error) {
	fmt.Printf("Mul function was invoked with: %v\n", req)
	value1 := req.GetMul().GetFirstNumber()
	value2 := req.GetMul().GetSecondNumber()
	result := value1 * value2
	res := &calculatorpb.MulResponse{
		MulResult: result,
	}
	return res, nil
}

func (*server) Div(ctx context.Context, req *calculatorpb.DivRequest) (*calculatorpb.DivResponse, error) {
	fmt.Printf("Div function was invoked with: %v\n", req)
	value1 := req.GetDiv().GetFirstNumber()
	value2 := req.GetDiv().GetSecondNumber()
	var result float32
	if value2 == 0 {
		result = 0.00
	} else {
		result = float32(value1 / value2)
	}
	res := &calculatorpb.DivResponse{
		DivResult: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello Simple Calculator!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
