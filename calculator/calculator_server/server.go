package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/neoyewchuan/go-grpc-course/calculator/calculatorpb"
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
	var result float64
	if value2 == 0 {
		result = 0.00
	} else {
		result = float64(value1) / float64(value2)
	}
	res := &calculatorpb.DivResponse{
		DivResult: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with: %v\n", req)
	number := req.GetNumber()
	primefactor := int64(2)
	for number > 1 {
		if number%primefactor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: primefactor,
			})
			number /= primefactor
		} else {
			primefactor++
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function was invoked with a streaming request..")
	counter := 0
	sumtotal := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have reading the client stream
			average := float64(sumtotal) / float64(counter)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}
		counter++
		sumtotal += req.GetNumber()

	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked with a streaming request..")
	var lastMax int64

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
			return err
		}
		theNumber := req.GetNumber()
		fmt.Printf("Received numer: %v\n", theNumber)
		if theNumber > lastMax {
			lastMax = theNumber
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				MaximumNumber: lastMax,
			})
			if err != nil {
				log.Fatalf("Error while sending data to client: %v\n", err)
				return err
			}
		}

	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number))
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil

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
