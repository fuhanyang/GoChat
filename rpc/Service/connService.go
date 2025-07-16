package Service

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"rpc/Client"
)

func ConnService(ServiceName string, builder resolver.Builder) (Client.Client, error) {

	//// 连接rpc服务
	//addr := GetServerAddr(ServiceName)
	//if addr == "" {
	//	return nil, fmt.Errorf("service not found")
	//}
	//fmt.Println(ServiceName, " service found at ", addr)
	etcdTarget := fmt.Sprintf("etcd:///%s", ServiceName)
	conn, err := grpc.Dial(fmt.Sprintf(etcdTarget),
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c := Client.NewClient(conn, ServiceName)
	fmt.Println("start listen ", ServiceName, " success")
	return c, nil
}
