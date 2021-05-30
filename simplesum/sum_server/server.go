package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/neoyewchuan/go-grpc-course/simplesum/sumpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with: %v\n", req)
	value1 := req.GetSum().GetValue_1()
	value2 := req.GetSum().GetValue_2()
	result := value1 + value2
	res := &sumpb.SumResponse{
		Result: result,
	}
	return res, nil
}
func main() {
	fmt.Println("Hello SimpleSum!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
