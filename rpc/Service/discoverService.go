package Service

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var (
	host = "127.0.0.1"
	port = "2379"
)

type Services struct {
	services map[string]*Service
	sync.RWMutex
}

var myServices = &Services{
	services: make(map[string]*Service),
}

func GetServerAddr(svcName string) string {
	s := ServiceDiscovery(svcName)
	if s == nil || (s.Host == "" && s.Port == "") {
		return ""
	}
	return s.Host + ":" + s.Port
}
func ServiceDiscovery(svcName string) *Service {
	myServices.RLock()
	defer myServices.RUnlock()
	var s *Service = nil
	s, _ = myServices.services[svcName]
	return s
}
func WatchServiceName(svcName string, host string, port string) error {
	fmt.Printf("start to watch service %s\n", svcName)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%s", host, port)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return err
	}
	defer cli.Close()

	Res, err := cli.Get(context.Background(), svcName, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if Res.Count > 0 {
		fmt.Printf("get service %s from etcd\n", svcName)
		mp := sliceToMap(Res.Kvs)
		s := &Service{}
		if kv, ok := mp[svcName]; ok {
			s.Name = string(kv.Value)
		}
		if kv, ok := mp[svcName+".ip"]; ok {
			s.Host = string(kv.Value)
		}
		if kv, ok := mp[svcName+".port"]; ok {
			s.Port = string(kv.Value)
		}
		if kv, ok := mp[svcName+".protocol"]; ok {
			s.Protocol = string(kv.Value)
		}
		// 获取到服务信息
		myServices.Lock()
		myServices.services[svcName] = s
		myServices.Unlock()
	}

	rch := cli.Watch(context.Background(), svcName, clientv3.WithPrefix())
	for wres := range rch {
		for _, ev := range wres.Events {
			if ev.Type == clientv3.EventTypeDelete {
				myServices.Lock()
				delete(myServices.services, svcName)
				myServices.Unlock()
			}
			// 新增或修改服务信息则更新本地缓存
			if ev.Type == clientv3.EventTypePut {
				myServices.Lock()
				if _, ok := myServices.services[svcName]; !ok {
					myServices.services[svcName] = &Service{}
				}
				switch string(ev.Kv.Key) {
				case svcName:
					myServices.services[svcName].Name = string(ev.Kv.Value)
				case svcName + ".ip":
					myServices.services[svcName].Host = string(ev.Kv.Value)
				case svcName + ".port":
					myServices.services[svcName].Port = string(ev.Kv.Value)
				case svcName + ".protocol":
					myServices.services[svcName].Protocol = string(ev.Kv.Value)

				}
				myServices.Unlock()
			}
		}
	}
	return nil
}

func sliceToMap(slice []*mvccpb.KeyValue) map[string]*mvccpb.KeyValue {
	mp := make(map[string]*mvccpb.KeyValue)
	for _, kv := range slice {
		mp[string(kv.Key)] = kv
	}
	return mp
}
