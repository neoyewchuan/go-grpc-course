package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/neoyewchuan/go-grpc-course/primenumberdecomposition/primenumdecompb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) PrimeDecomManyTimes(req *primedecompb.PrimeDecomManyTimesRequest, stream primedecompb.PrimeDecomService_PrimeDecomManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with: %v\n", req)
	number := req.GetPrimedecom().GetNumber()
	k := 2
	for number > 1 {
		if number%k == 0 {
			res := &primedecompb.PrimeDecomManyTimesResponse{
				Result: k,
			}
			number = number / k
			stream.Send(res)
		} else {
			k = k + 1
		}
	}

	return nil
}

func main() {
	fmt.Println("Prime Number Decomposition")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	primedecompb.RegisterPrimeDesomServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
