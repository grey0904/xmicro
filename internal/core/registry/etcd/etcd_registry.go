package etcd

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
	"xmicro/internal/utils/u_conv"
	// "go.etcd.io/etcd/client/v3" // 假设你在使用Etcd的Go客户端
)

type Registry struct {
	cli         *clientv3.Client                        //etcd连接
	leaseId     clientv3.LeaseID                        //租约id
	DialTimeout int                                     //超时时间
	ttl         int                                     //租约时间
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse //心跳
	info        Server                                  //注册的server信息
	closeCh     chan struct{}
}

var once sync.Once
var reg *Registry

func NewEtcdRegistry() *Registry {
	once.Do(func() {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   config.Conf.Etcd.Addrs,
			DialTimeout: time.Duration(config.Conf.Etcd.DialTimeout) * time.Second,
		})
		if err != nil {
			logs.Error("clientv3.New err:%v", err)
			return
		}
		reg = &Registry{
			cli:         client,
			DialTimeout: 3,
		}
	})
	return reg
}

func (r *Registry) Register() error {
	r.info = Server{
		Name:    config.LocalConf.AppName,
		Addr:    config.Conf.ServerRpc.Host + ":" + u_conv.Uint64ToString(config.Conf.ServerRpc.Port),
		Weight:  config.Conf.Etcd.Register.Weight,
		Version: config.Conf.Etcd.Register.Version,
		Ttl:     config.Conf.Etcd.Register.Ttl,
	}

	if err := r.register(); err != nil {
		return err
	}
	r.closeCh = make(chan struct{})

	//放入协程中 根据心跳的结果 做相应的操作
	go r.watcher()
	return nil
}

func (r *Registry) register() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(r.DialTimeout))
	defer cancel()

	//1. 创建租约
	grant, err := r.cli.Grant(ctx, r.info.Ttl)
	if err != nil {
		logs.Error("createLease failed,err:%v", err)
		return err
	}
	r.leaseId = grant.ID

	//2. 心跳检测
	//心跳 要求是一个长连接 如果做了超时 长连接就断掉了 不要设置超时
	//就是一直不停的发消息 保持租约 续租
	r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leaseId)
	if err != nil {
		logs.Error("keepAlive failed,err:%v", err)
		return err
	}

	//3. 绑定租约
	data, _ := json.Marshal(r.info)
	key := r.info.BuildRegisterKey()
	_, err = r.cli.Put(ctx, key, string(data), clientv3.WithLease(r.leaseId))
	if err != nil {
		logs.Error("bindLease failed,err:%v", err)
		return err
	}

	logs.Info("register service success,key=%s", key)
	return nil
}

// watcher 续约 新注册 close 注销
func (r *Registry) watcher() {
	//租约到期了 是不是需要去检查是否自动注册
	ticker := time.NewTicker(time.Duration(r.info.Ttl) * time.Second)
	for {
		select {
		case res := <-r.keepAliveCh:
			//如果etcd重启了 相当于连接断开 需要进行重新连接 res==nil
			if res == nil {
				if err := r.register(); err != nil {
					logs.Error("keepAliveCh register failed,err:%v", err)
				}
				logs.Info("续约重新注册成功,%v", res)
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					logs.Error("ticker register failed,err:%v", err)
				}
			}
		}
	}
}

func (r *Registry) Deregister() error {
	_, err := r.cli.Delete(context.Background(), r.info.BuildRegisterKey())
	if err != nil {
		logs.Error("close and unregister failed,err:%v", err)
		return err
	}
	//租约撤销
	if _, err = r.cli.Revoke(context.Background(), r.leaseId); err != nil {
		logs.Error("close and Revoke lease failed,err:%v", err)
		return err
	}
	if r.cli != nil {
		r.cli.Close()
	}
	logs.Info("unregister etcd...")

	return nil
}

//func (r *Registry) Discover() ([]registry.InstanceInfo, error) {
//	// Nacos服务发现逻辑
//	return nil, nil
//}
