package main

import (
	"context"
	"fmt"
	dtproto "github.com/aaronchen2k/deeptest/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
)

func main() {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := dtproto.NewTestServiceClient(connect)

	stream, err := client.TestBidirectional(context.Background())
	if err != nil {
		log.Println(err.Error())
		return
	}

	// send
	go func() {
		for true {
			err = stream.Send(&dtproto.TestData{
				Username: "aaron",
				Email:    "462826@qq.com",
			})
		}

		err := stream.CloseSend()
		if err != nil {
			log.Println(err)
		}
	}()

	// recv
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for true {
			resp, err := stream.Recv()
			//if err == io.EOF {
			//	break
			//}

			log.Println(fmt.Sprintf("get msg from grpc producer %v", resp))

			if err != nil {
				continue
			}
		}

		stream.CloseSend()
	}()

	wg.Wait()
}
