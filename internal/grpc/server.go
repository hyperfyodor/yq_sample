package grpc

import (
	"context"
	service "github.com/hyperfyodor/yq_sample/internal/service/consumer"
	"github.com/hyperfyodor/yq_sample/proto/consumer/gen"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConsumerServer struct {
	gen.UnimplementedConsumerServiceServer
	consumerService *service.Consumer
	limiter         *rate.Limiter
}

func NewConsumerServer(consumer *service.Consumer, limiter *rate.Limiter) *ConsumerServer {
	if consumer == nil {
		panic("consumer service is nil")
	}
	if limiter == nil {
		panic("limiter is nil")
	}

	return &ConsumerServer{consumerService: consumer, limiter: limiter}
}

func (server *ConsumerServer) ProcessTask(ctx context.Context, req *gen.ProcessTaskRequest) (*gen.ProcessTaskResponse, error) {
	if err := server.limiter.Wait(ctx); err != nil {
		return &gen.ProcessTaskResponse{}, status.Error(codes.Internal, "failed while waiting limiter to open")
	}

	if err := server.consumerService.Consume(ctx, int(req.Id), int(req.Type), int(req.Value)); err != nil {
		return &gen.ProcessTaskResponse{}, status.Error(codes.Internal, "failed while consuming task")
	}

	return &gen.ProcessTaskResponse{}, nil
}
