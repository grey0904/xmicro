package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"xmicro/internal/app"
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

func InitRedis() {
	c := app.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password, // no password set
		DB:       c.Select,   // use default DB
		//PoolSize:     15,
		//MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("initMysqlConfig yaml.Unmarshal err: %v", err)
	}
	app.Rd = rdb
}
