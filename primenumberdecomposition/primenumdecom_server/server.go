package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/neoyewchuan/go-grpc-course/primenumberdecomposition/primenumdecompb"
)

type server struct{}

func (*server) PrimeDecomManyTimes(req *primenumdecompb.PrimeDecomManyTimesRequest, stream primenumdecompb.PrimeDecomService_PrimeDecomManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with: %v\n", req)
	number := req.GetPrimedecom().GetNumber()
	var k int32
	k = 2
	for number > 1 {
		if number%k == 0 {
			res := &primenumdecompb.PrimeDecomManyTimesResponse{
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
	primenumdecompb.RegisterPrimeDecomServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
