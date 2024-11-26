package grpc

import (
	"context"
	"github.com/hyperfyodor/yq_sample/internal/helpers"
	"github.com/hyperfyodor/yq_sample/proto/consumer/gen"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type GrpcPublisher struct {
	client  gen.ConsumerServiceClient
	limiter *rate.Limiter
}

func NewGrpcPublisher(host string, port string, mps int) (*GrpcPublisher, error) {
	address := net.JoinHostPort(host, port)
	cc, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := gen.NewConsumerServiceClient(cc)

	limiter := rate.NewLimiter(rate.Limit(mps), 1)

	return &GrpcPublisher{client, limiter}, nil
}

// Publish implementation of the interface, instead of using tick to produce N messages per second,
// we will produce the maximum possible, rate limited with mps
func (publisher *GrpcPublisher) Publish(ctx context.Context, taskId int, taskType int, taskValue int) error {
	const op = "internal.grpc.GrpcPublisher.Publish"

	select {
	case <-ctx.Done():
		return helpers.WrapErr(op, ctx.Err())
	default:
		if err := publisher.limiter.Wait(ctx); err != nil {
			return helpers.WrapErr(op, err)
		}
		_, err := publisher.client.ProcessTask(
			ctx,
			&gen.ProcessTaskRequest{Id: int32(taskId), Type: int32(taskType), Value: int32(taskValue)},
			grpc.WaitForReady(false),
		)

		if err != nil {
			return err
		}

		return nil
	}
}
