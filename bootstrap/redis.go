package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"xmicro/internal/app"
	"xmicro/internal/log"
)

var ctx = context.Background()

// redis 集群
//func InitRedisCluster() *redis.ClusterClient {
//	rdcInfo := config.Conf.RedisCluster
//	rdb := redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs: []string{rdcInfo[0], rdcInfo[1], rdcInfo[2], rdcInfo[3], rdcInfo[4], rdcInfo[5]},
//	})
//	err := rdb.Ping(ctx).Err()
//	if err != nil {
//		return nil
//	}
//	return rdb
//}

type Redis struct {
	Host     string `json:"user" yaml:"host"`
	Password string `json:"password" yaml:"password"`
	Port     int    `json:"port" yaml:"port"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
	Select   int    `json:"select" yaml:"select"`
}

func InitRedis() error {
	content, err := app.App.NacosClient.GetConfig(vo.ConfigParam{
		DataId: "grey_redis.yaml",
	})
	if err != nil {
		log.Logger.Error("nacos GetConfig 错误:", err)
		return err
	}

	var r Redis

	err = yaml.Unmarshal([]byte(content), &r)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(content), &r)
	if err != nil {
		log.Logger.Error("Redis 初始化错误:", err)
		return err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
		Password: r.Password, // no password set
		DB:       r.Select,   // use default DB
		//PoolSize:     15,
		//MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	})
	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Logger.Error("Redis 初始化错误:", err)
		return err
	}

	app.App.Redis = rdb
	return nil
}
