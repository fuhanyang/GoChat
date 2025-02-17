package UserPb

import (
	"google.golang.org/grpc"
	"rpc/User"
)

func NewClient(conn *grpc.ClientConn) User.UserServiceClient {
	c := User.NewUserServiceClient(conn)
	return c
}
