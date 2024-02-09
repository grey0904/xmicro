package app

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Redis       *redis.Client
	DB          *gorm.DB
	NacosClient config_client.IConfigClient
)
