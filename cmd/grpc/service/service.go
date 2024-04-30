package service

import (
	"context"
	dtproto "github.com/aaronchen2k/deeptest/proto"
	"io"
	"log"
	"time"
)

type GrpcService struct {
	execCtx    context.Context
	execCancel context.CancelFunc
}

func (s *GrpcService) TestBidirectional(stream dtproto.TestService_TestBidirectionalServer) (err error) {
	// get request
	req, err := stream.Recv()
	if err == io.EOF {
		err = nil
	}
	if req == nil {
		return
	}

	s.execCtx, _ = context.WithTimeout(context.Background(), time.Duration(6)*time.Second)

	go func() {
		log.Println(">>>>>> start")

		for true {
			time.Sleep(1 * time.Second)

			log.Println("====== do something")

			select {
			case <-s.execCtx.Done():
				goto LabelEnd

			default:
			}
		}

	LabelEnd:
		log.Println("<<<<<<< end")
	}()

	return
}

func (s *GrpcService) TestUnidirectional(ctx context.Context, data *dtproto.TestData) (ret *dtproto.TestResponse, err error) {
	ret.Code = 0
	ret.Data = data

	return
}
