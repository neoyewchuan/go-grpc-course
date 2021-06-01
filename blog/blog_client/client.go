package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/neoyewchuan/go-grpc-course/blog/blogpb"
)

func main() {

	fmt.Println("Blog Client connecting...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	fmt.Println("Blog Client running...")
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	//CreateBlog(c)
	ReadBlog(c)

}

func CreateBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Creating blog...")
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "neo",
			Title:    "Genesis Blog",
			Content:  "My very first blog to be created in the gRPC CRUD API on MongoDB",
		},
	}
	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Blog RPC: %v", err)
	}
	log.Printf("Blog has been created with Id: %v", res)
}

func ReadBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Reading blog...")
	req := &blogpb.ReadBlogRequest{
		BlogId: "60b65d063fb6b015585ba906",
	}
	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Blog RPC: %v", err)
	}
	log.Printf("Blog has been retrieved: %v", res)
}
