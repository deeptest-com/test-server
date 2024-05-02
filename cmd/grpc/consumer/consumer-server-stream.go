package main

import (
	"context"
	"fmt"
	dtproto "github.com/aaronchen2k/deeptest/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"sync"
)

func main() {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := dtproto.NewTestServiceClient(connect)

	req := dtproto.TestRequest{
		Action: "act",
		Data:   "=== server-stream",
	}

	stream, err := client.TestServerStream(context.Background(), &req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for true {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}

			log.Println(fmt.Sprintf("get msg from grpc producer %v", resp))

			if err != nil {
				continue
			}
		}
	}()

	wg.Wait()
}
