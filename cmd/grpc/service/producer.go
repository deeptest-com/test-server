package service

import (
	"context"
	dtproto "github.com/aaronchen2k/deeptest/proto"
	"io"
	"log"
)

type GrpcService struct {
	execCtx    context.Context
	execCancel context.CancelFunc
}

func (s *GrpcService) TestBidirectional(stream dtproto.TestService_TestBidirectionalServer) (err error) {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("end of stream")
			break
		}

		log.Printf("got msg from grpc client %v", req)

		resp := dtproto.TestResponse{
			Code: 456,
			Data: req,
		}

		err = stream.Send(&resp)
		if err != nil {
			log.Printf("send msg error %v", err)
		}
	}

	return
}

func (s *GrpcService) TestUnidirectional(ctx context.Context, data *dtproto.TestData) (ret *dtproto.TestResponse, err error) {
	ret = &dtproto.TestResponse{
		Code: 123,
		Data: data,
	}

	return
}
