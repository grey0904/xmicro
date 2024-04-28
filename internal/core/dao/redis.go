package dao

import (
	"context"
	"fmt"
	"xmicro/internal/core/repo"
)

const Prefix = "MSQP"
const AccountIdRedisKey = "AccountId"
const AccountIdBegin = 10000

type RedisDao struct {
	repo *repo.Manager
}

func (d *RedisDao) NextAccountId() (string, error) {
	//自增 给一个前缀
	return d.incr(Prefix + ":" + AccountIdRedisKey)
}

func (d *RedisDao) incr(key string) (string, error) {
	//判断此key是否存在 不存在 set 存在就自增
	todo := context.TODO()
	var exist int64
	var err error
	//0 代表不存在
	if d.repo.Redis.Cli != nil {
		exist, err = d.repo.Redis.Cli.Exists(todo, key).Result()
	} else {
		exist, err = d.repo.Redis.ClusterCli.Exists(todo, key).Result()
	}
	if exist == 0 {
		//不存在
		if d.repo.Redis.Cli != nil {
			err = d.repo.Redis.Cli.Set(todo, key, AccountIdBegin, 0).Err()
		} else {
			err = d.repo.Redis.ClusterCli.Set(todo, key, AccountIdBegin, 0).Err()
		}
		if err != nil {
			return "", err
		}
	}
	var id int64
	if d.repo.Redis.Cli != nil {
		id, err = d.repo.Redis.Cli.Incr(todo, key).Result()
	} else {
		id, err = d.repo.Redis.ClusterCli.Incr(todo, key).Result()
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}

func NewRedisDao(m *repo.Manager) *RedisDao {
	return &RedisDao{
		repo: m,
	}
}
