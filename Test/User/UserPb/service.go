package UserPb

import (
	"google.golang.org/grpc"
	"rpc/user"
)

func NewClient(conn *grpc.ClientConn) user.UserServiceClient {
	c := user.NewUserServiceClient(conn)
	return c
}
