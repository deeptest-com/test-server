package producer

import (
	"context"
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

			log.Printf("got msg from grpc client %v", req)

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
		}
	}

	return
}

func (s *GrpcService) TestClientStream(server dtproto.TestService_TestClientStreamServer) (err error) {
	for {
		req, err := server.Recv()
		if err == io.EOF {
			log.Println("end of stream")
			break
		}

		log.Printf("got msg from grpc client %v", req)
	}

	resp := dtproto.TestResponse{
		Code:   456,
		Result: "success",
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
