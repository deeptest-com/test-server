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

	stream, err := client.TestClientStream(context.Background())
	if err != nil {
		log.Println(err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			act := "do"
			if i == 3 {
				act = "stop"
			}

			err = stream.Send(&dtproto.TestRequest{
				Action: act,
				Data:   "=== client-stream",
			})

			if err != nil {
				log.Printf("send msg error %v", err)
				continue
			}
		}

		stream.CloseSend()
	}()

	wg.Wait()

	resp, err := stream.CloseAndRecv()
	log.Println(fmt.Sprintf("get msg from grpc producer %v", resp))
}
