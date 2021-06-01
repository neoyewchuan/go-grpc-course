package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"github.com/neoyewchuan/go-grpc-course/greet/greetpb"
)

// ignoreCN disables interpreting Common Name as a hostname. See issue 24151.
var ignoreCN = !strings.Contains(os.Getenv("GODEBUG"), "x509ignoreCN=0")

func main() {
	fmt.Println("Hello I'm a client...")

	tls := true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "./ssl/ca.crt" // CA trust certificate

		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}
	cc, err := grpc.Dial("localhost:50051", opts)
	//cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client: %f", c)

	doUnary(c)
	doServerStreaming(c)
	doClientStreaming(c)
	doBiDiStreaming(c)

	doUnaryWithDeadline(c, &greetpb.Greeting{FirstName: "Peter", LastName: "Parker"}, 1*time.Second) // should timeout
	doUnaryWithDeadline(c, &greetpb.Greeting{FirstName: "Harry", LastName: "Potter"}, 5*time.Second) // should complete
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Peter",
			LastName:  "Parker",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephanie",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Angelina",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Cynthia",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Michelle",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Zune",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Gracelyn",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	log.Printf("LongGreet Response: %v\n", res)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
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
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v\n", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephanie",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Angelina",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Cynthia",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Michelle",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Zune",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Gracelyn",
			},
		},
	}

	waitc := make(chan struct{})

	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
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
			fmt.Printf("Received message: %v\n", res.GetResult())
			time.Sleep(1000 * time.Millisecond)
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, greeting *greetpb.Greeting, timeout time.Duration) {
	fmt.Println("Starting to do Unary with Deadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: greeting,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hot! Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline RPC: %v\n", err)
		}
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
