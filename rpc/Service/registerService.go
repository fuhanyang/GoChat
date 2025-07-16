package Service

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type Service struct {
	Name     string
	Host     string
	Port     string
	Protocol string
}

// ServiceRegister 注册服务
func ServiceRegister(s *Service, ctx context.Context, host string, port string) error {
	// 构造服务地址
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	//连接etcd
	etcdUrl := fmt.Sprintf("http://%s:%s", host, port)
	cli, err := clientv3.NewFromURL(etcdUrl)
	if err != nil {
		return err
	}
	etcdManager, err := endpoints.NewManager(cli, s.Name)
	if err != nil {
		return err
	}
	//cli, err := clientv3.New(clientv3.Config{
	//	Endpoints:   []string{dsn},
	//	DialTimeout: 5 * time.Second,
	//})
	//if err != nil {
	//	fmt.Printf("connect to etcd failed, err:%v\n", err)
	//	return err
	//}

	var grantLease bool
	var leaseID clientv3.LeaseID

	res, err := cli.Get(ctx, s.Name, clientv3.WithCountOnly())
	if err != nil {
		return err
	}
	if res.Count == 0 {
		// 需要分配租约
		grantLease = true
	} else {
		fmt.Println("service already exists")
	}
	if grantLease {
		leaseRes, err := cli.Grant(ctx, 10)
		if err != nil {
			return err
		}
		leaseID = leaseRes.ID
		fmt.Printf("lease id = %v\n", leaseID)
	}

	// 注册服务
	err = etcdManager.AddEndpoint(ctx, fmt.Sprintf("%s/%s", s.Name, addr), endpoints.Endpoint{Addr: addr}, clientv3.WithLease(leaseID))
	if err != nil {
		return err
	}
	//kv := clientv3.NewKV(cli)
	//txn := kv.Txn(ctx)
	//// 判断key是否存在，不存在则创建，存在则更新
	//_, err = txn.If(clientv3.Compare(clientv3.CreateRevision(s.Name), "=", 0)).
	//	Then(
	//		clientv3.OpPut(s.Name, s.Name, clientv3.WithLease(leaseID)),
	//		clientv3.OpPut(s.Name+".ip", s.Host, clientv3.WithLease(leaseID)),
	//		clientv3.OpPut(s.Name+".port", s.Port, clientv3.WithLease(leaseID)),
	//		clientv3.OpPut(s.Name+".protocol", s.Protocol, clientv3.WithLease(leaseID)),
	//	).
	//	Else(
	//		clientv3.OpPut(s.Name, s.Name, clientv3.WithIgnoreLease()),
	//		clientv3.OpPut(s.Name+".ip", s.Host, clientv3.WithIgnoreLease()),
	//		clientv3.OpPut(s.Name+".port", s.Port, clientv3.WithIgnoreLease()),
	//		clientv3.OpPut(s.Name+".protocol", s.Protocol, clientv3.WithIgnoreLease()),
	//	).
	//	Commit()
	//if err != nil {
	//	return err
	//}

	go func() {
		defer cli.Close()
		if grantLease {
			// 续租
			ctx := context.Background()
			leaseKeepAlive, err := cli.KeepAlive(ctx, leaseID)
			if err != nil {
				fmt.Printf("keep alive failed, err:%v\n", err)
				return
			}
			for {
				select {
				case lease := <-leaseKeepAlive:
					if lease == nil {
						fmt.Printf("lease keep alive channel closed\n")
						return
					}
				case <-ctx.Done():
				}
			}
		}
	}()
	return nil
}
