package main

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/cmd/grpc/producer"
	dtproto "github.com/aaronchen2k/deeptest/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = "0.0.0.0:9528"
)

func main() {
	startGrpcServe()
}

func startGrpcServe() {
	server := grpc.NewServer()
	dtproto.RegisterTestServiceServer(server, &producer.GrpcService{})
	reflection.Register(server)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("grpc net.Listen err: %v", err)
	}

	log.Println(fmt.Sprintf("start grpc on %s", port))

	server.Serve(lis)
}
