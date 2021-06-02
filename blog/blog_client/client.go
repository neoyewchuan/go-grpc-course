package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"

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
	//ReadBlog(c)
	//UpdateBlog(c)
	//DeleteBlog(c)
	ListBlog(c)

}

func CreateBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Create a blog...")
	author := []string{"lina", "shaoqi", "doreen", "cynthia", "michelle", "gracelyn", "neo"}
	title := []string{"happy moments", "happy birthday", "forget me not", "learning programming", "google golang & grpc", "no title", "my favourite food"}
	content := []string{"some content", "some random content", "more random content", "don't know what to write", "I'm speechless", "I'm jobless, please hire me", "fried carrot cake, nasi lemak, roti prata, economic bee hoon"}

	for i := 0; i < 3; i++ {
		idx := rand.Intn(7)

		req := &blogpb.CreateBlogRequest{
			Blog: &blogpb.Blog{
				AuthorId: author[idx],
				Title:    title[idx],
				Content:  content[idx],
			},
		}
		res, err := c.CreateBlog(context.Background(), req)
		if err != nil {
			log.Fatalf("error while calling Blog RPC: %v", err)
		}
		log.Printf("Blog has been created with Id: %v\n", res)
	}

}

func ReadBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Read a blog...")
	req := &blogpb.ReadBlogRequest{
		BlogId: "60b65d063fb6b015585ba906",
	}
	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Blog RPC: %v", err)
	}
	log.Printf("Blog has been retrieved: %v\n", res)
}

func UpdateBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Update a blog...")

	newBlock := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "60b707e287ae1dcdb9c6dd23",
			AuthorId: "shaoqi",
			Title:    "Blog Star",
			Content:  "This is an update blog for previously created blog {60b65d063fb6b015585ba906}",
		},
	}
	res, err := c.UpdateBlog(context.Background(), newBlock)
	if err != nil {
		log.Fatalf("error while updating Blog: %v", err)
	}
	fmt.Printf("Blog has been updated: %v\n", res)
}

func DeleteBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Delete a blog...")

	blogToDel := "60b7730e40ae4ca78bd639b6"

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogToDel})
	if err != nil {
		log.Fatalf("error while calling Blog RPC: %v", err)
	}
	if res != nil {
		req := &blogpb.DeleteBlogRequest{
			BlogId: blogToDel,
		}
		delRes, delErr := c.DeleteBlog(context.Background(), req)
		if delErr != nil {
			log.Fatalf("error while deleting Blog: %v", delErr)
		}
		fmt.Printf("Blog has been deleted: %v\n", delRes)

	}

}

func ListBlog(c blogpb.BlogServiceClient) {
	fmt.Println("List Blog entries...")
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		fmt.Printf("Error while calling ListBlog Request: %v", err)
	}
	if stream != nil {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Something happened: %v", err)
			}
			fmt.Printf("A Blog Entry Discovered: %v\n", res.GetBlog())
			//time.Sleep(500 * time.Millisecond)
		}
	} else {
		fmt.Println("Nothing in the stream, no blog entry found!")
	}

}
