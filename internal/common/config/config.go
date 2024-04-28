package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"strconv"
	"strings"
)

var LocalConf *LocalConfig
var Conf *Config
var Nc config_client.IConfigClient

const (
	MysqlConfigKey = "mysql.yaml"
	RedisConfigKey = "redis.yaml"
	MongoConfigKey = "mongo.yaml"
)

// InitConfig 加载配置
func InitConfig(appName string) {
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configFile == "" {
		panic(fmt.Errorf("path to config file err"))
	}

	LocalConf = new(LocalConfig)
	v := viper.New()
	v.SetConfigFile(*configFile)
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件被修改了")
		if err := v.Unmarshal(&LocalConf); err != nil {
			panic(fmt.Errorf("Unmarshal change config data,err:%v \n", err))
		}
	})

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件出错,err:%v \n", err))
	}

	if err := v.Unmarshal(&LocalConf); err != nil {
		panic(fmt.Errorf("Unmarshal config data,err:%v \n", err))
	}

	LocalConf.AppName = appName

	// 用 AppConfig 中的Nacos配置信息创建“配置中心客户端”
	newConfigClient()
	// 从Nacos上获取配置，并解析给对应的结构体
	initAppConfig()
	initDatabaseConfigs()
}

func newConfigClient() {
	var (
		sc = make([]constant.ServerConfig, 0)
		nc = LocalConf.Nacos
	)

	cc := constant.ClientConfig{
		Username:            nc.Username,
		Password:            nc.Password,
		TimeoutMs:           nc.TimeoutMs,
		NamespaceId:         nc.NamespaceId, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		NotLoadCacheAtStart: nc.NotLoadCacheAtStart,
		CacheDir:            nc.CacheDir,
		LogDir:              nc.LogDir,
		LogLevel:            nc.LogLevel,
	}

	for _, value := range nc.Endpoints {
		vs := strings.Split(value, ":")
		if len(vs) < 2 {
			log.Fatalf("创建Nacos客户端失败:endpoints 配置有误")
		}

		port, err := strconv.ParseUint(vs[1], 10, 64)
		if err != nil {
			log.Fatalf("创建Nacos客户端失败:endpoints 配置有误")
		}

		sc = append(sc, constant.ServerConfig{
			IpAddr: vs[0],
			Port:   port,
		})
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		log.Fatalf("创建Nacos客户端失败:endpoints 配置有误, %s", err)
	}

	Nc = client
}

func initAppConfig() {
	content, err := Nc.GetConfig(vo.ConfigParam{
		DataId: LocalConf.AppName + ".yaml",
	})
	if err != nil {
		log.Fatalf("InitAppConfig NacosClient.GetConfig err: %v", err)
	}

	err = yaml.Unmarshal([]byte(content), &Conf)
	if err != nil {
		log.Fatalf("InitAppConfig yaml.Unmarshal err: %v", err)
	}
}

func initDatabaseConfigs() {
	initConfig(MysqlConfigKey, &Conf.Database.MysqlConf)
	initConfig(RedisConfigKey, &Conf.Database.RedisConf)
	initConfig(MongoConfigKey, &Conf.Database.MongoConf)
}

func initConfig(configKey string, target interface{}) {
	content, err := Nc.GetConfig(vo.ConfigParam{
		DataId: configKey,
	})
	if err != nil {
		log.Fatalf("Failed to initialize %s config: %v", configKey, err)
	}

	err = yaml.Unmarshal([]byte(content), target)
	if err != nil {
		log.Fatalf("Failed to unmarshal %s config: %v", configKey, err)
	}
}
