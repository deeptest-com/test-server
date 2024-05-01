package main

import (
	"context"
	"fmt"
	dtproto "github.com/aaronchen2k/deeptest/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := dtproto.NewTestServiceClient(connect)

	req := dtproto.TestRequest{
		Action: "act",
		Data:   "===",
	}

	resp, err := client.TestUnary(context.Background(), &req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(fmt.Sprintf("get msg from grpc producer %v", resp))
}
