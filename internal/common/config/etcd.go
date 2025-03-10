package config

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdCenter struct {
	client *clientv3.Client
}

func NewEtcdCenter(cfg interface{}) (*EtcdCenter, error) {
	etcdConfig, ok := cfg.(EtcdClientConfig)
	if !ok {
		return nil, fmt.Errorf("invalid etcd config type")
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdConfig.Endpoints,
		Username:    etcdConfig.Username,
		Password:    etcdConfig.Password,
		DialTimeout: etcdConfig.DialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %v", err)
	}

	return &EtcdCenter{client: client}, nil
}

func (e *EtcdCenter) GetConfig(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get config from etcd: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("config not found: %s", key)
	}

	return string(resp.Kvs[0].Value), nil
}

func (e *EtcdCenter) WatchConfig(key string, onChange func(string)) error {
	go func() {
		watchChan := e.client.Watch(context.Background(), key)
		for resp := range watchChan {
			for _, ev := range resp.Events {
				if ev.Type == clientv3.EventTypePut {
					onChange(string(ev.Kv.Value))
				}
			}
		}
	}()
	return nil
}

func (e *EtcdCenter) Close() error {
	return e.client.Close()
}
