package Client

import (
	"context"
	"google.golang.org/grpc"
	"rpc/handler"
)

type ClientFactory func(conn *grpc.ClientConn) Client

var clientFactoryMap = map[string]ClientFactory{}

func Inject(serviceName string, factory ClientFactory) {
	clientFactoryMap[serviceName] = factory
}

type Client interface {
	Handle(ctx context.Context, Handler handler.HandlerRequest) ([]byte, error)
	InjectClientFactory()
}

func NewClient(conn *grpc.ClientConn, ServiceName string) Client {
	return clientFactoryMap[ServiceName](conn)
}
