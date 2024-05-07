package producer

import (
	"context"
	"fmt"
	dtproto "github.com/aaronchen2k/deeptest/proto"
	"io"
	"log"
	"sync"
)

type GrpcService struct {
	execCtx    context.Context
	execCancel context.CancelFunc
}

func (s *GrpcService) TestBidi(server dtproto.TestService_TestBidiServer) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)

	finish := false

	go func() {
		for {
			req, err := server.Recv()
			if err == io.EOF {
				log.Println("end of stream")
				//break
			}

			log.Printf("got msg from grpc consumer %v", req)

			if req.Action == "stop" {
				finish = true
				wg.Done()
				break
			}
		}
	}()

	go func() {
		for true {
			resp := dtproto.TestResponse{
				Code:   456,
				Result: "success",
			}

			err = server.Send(&resp)
			if err != nil {
				log.Printf("send msg error %v", err)
			}

			if finish {
				break
			}
		}
	}()

	wg.Wait()

	return
}

func (s *GrpcService) TestServerStream(req *dtproto.TestRequest, server dtproto.TestService_TestServerStreamServer) (err error) {
	log.Printf("got msg from grpc client %v", req)

	for {
		resp := dtproto.TestResponse{
			Code:   456,
			Result: "success",
		}

		err = server.Send(&resp)
		if err != nil {
			log.Printf("send msg error %v", err)

			break
		}
	}

	return
}

func (s *GrpcService) TestClientStream(server dtproto.TestService_TestClientStreamServer) (err error) {
	count := 0
	for i := 0; i < 10; i++ {
		req, err := server.Recv()
		if err == io.EOF {
			log.Println("end of stream")
			break
		}
		if err != nil {
			log.Println(err.Error())
			break
		}

		log.Printf("got msg from grpc client %v", req)
		count++

		if req.Action == "stop" {
			break
		}
	}

	resp := dtproto.TestResponse{
		Code:   456,
		Result: fmt.Sprintf("success, got %d msg", count),
	}

	err = server.SendAndClose(&resp)
	if err != nil {
		log.Printf("send msg error %v", err)
	}

	return
}

func (s *GrpcService) TestUnary(ctx context.Context, request *dtproto.TestRequest) (resp *dtproto.TestResponse, err error) {
	resp = &dtproto.TestResponse{
		Code:   456,
		Result: "success",
	}

	return
}
