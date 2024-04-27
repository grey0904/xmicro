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
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"log"
	"strconv"
	"strings"
)

var NacosConf *NacosConfig
var Conf *Config
var Nc config_client.IConfigClient

// NacosConfig 本地的 nacos 配置
type NacosConfig struct {
	Nacos Nacos `yaml:"nacos"`
}

type Nacos struct {
	Endpoints           []string `yaml:"endpoints"`
	Username            string   `yaml:"username"`
	Password            string   `yaml:"password"`
	TimeoutMs           uint64   `yaml:"timeoutMs"`
	NamespaceId         string   `yaml:"namespaceId"`
	NotLoadCacheAtStart bool     `yaml:"otLoadCacheAtStart"`
	CacheDir            string   `yaml:"cacheDir"`
	LogDir              string   `yaml:"logDir"`
	LogLevel            string   `yaml:"logLevel"`
}

// Config 服务相关配置，配置存放在
type Config struct {
	Database   Database   `yaml:"db"`
	MetricPort int        `yaml:"metricPort"`
	AppName    string     `yaml:"appName"`
	ZapLog     ZapLogConf `yaml:"zapLog"`
	Etcd       EtcdConf   `yaml:"etcd"`
	Grpc       GrpcConf   `yaml:"grpc"`
}

type RegisterServer struct {
	Addr    string `yaml:"addr"`
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Weight  int    `yaml:"weight"`
	Ttl     int64  `yaml:"ttl"` //租约时长
}

// Database 数据库配置
type Database struct {
	MysqlConf MysqlConf `yaml:"mysql"`
	RedisConf RedisConf `yaml:"redis"`
	MongoConf MongoConf `yaml:"mongo"`
}

type MysqlConf struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Database     string `yaml:"database"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
}

type RedisConf struct {
	Addr         string   `yaml:"addr"`
	ClusterAddrs []string `yaml:"clusterAddrs"`
	Password     string   `yaml:"password"`
	PoolSize     int      `yaml:"poolSize"`
	MinIdleConns int      `yaml:"minIdleConns"`
	Host         string   `yaml:"host"`
	Port         int      `yaml:"port"`
	Timeout      int      `json:"timeout" yaml:"timeout"`
	Select       int      `json:"select" yaml:"select"`
}

type MongoConf struct {
	Url         string `yaml:"url"`
	Db          string `yaml:"db"`
	UserName    string `yaml:"userName"`
	Password    string `yaml:"password"`
	MinPoolSize int    `yaml:"minPoolSize"`
	MaxPoolSize int    `yaml:"maxPoolSize"`
}

type ZapLogConf struct {
	Level zapcore.Level `yaml:"level"`
	File  string        `yaml:"file"`
}

type EtcdConf struct {
	Addrs       []string       `yaml:"addrs"`
	RWTimeout   int            `yaml:"rwTimeout"`
	DialTimeout int            `yaml:"dialTimeout"`
	Register    RegisterServer `yaml:"register"`
}

type GrpcConf struct {
	Addr string `yaml:"addr"`
}

// InitConfig 加载配置
func InitConfig() {
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configFile == "" {
		panic(fmt.Errorf("path to config file err"))
	}

	NacosConf = new(NacosConfig)
	v := viper.New()
	v.SetConfigFile(*configFile)
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件被修改了")
		if err := v.Unmarshal(&NacosConf); err != nil {
			panic(fmt.Errorf("Unmarshal change config data,err:%v \n", err))
		}
	})

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件出错,err:%v \n", err))
	}

	if err := v.Unmarshal(&NacosConf); err != nil {
		panic(fmt.Errorf("Unmarshal config data,err:%v \n", err))
	}

	// 用 AppConfig 中的Nacos配置信息创建“配置中心客户端”
	NewConfigClient()
	// 从Nacos上获取Mysql、Redis等配置，并解析给对应的 AppConfig 里面的结构体
	InitMysqlConfig()
	InitRedisConfig()
}

func NewConfigClient() {
	var (
		sc = make([]constant.ServerConfig, 0)
		nc = NacosConf.Nacos
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

func InitMysqlConfig() {
	content, err := Nc.GetConfig(vo.ConfigParam{
		DataId: "mysql.yaml",
	})
	if err != nil {
		log.Fatalf("initMysqlConfig NacosClient.GetConfig err: %v", err)
	}

	err = yaml.Unmarshal([]byte(content), &Conf.Database.MysqlConf)
	if err != nil {
		log.Fatalf("initMysqlConfig yaml.Unmarshal err: %v", err)
	}
}

func InitRedisConfig() {
	content, err := Nc.GetConfig(vo.ConfigParam{
		DataId: "redis.yaml",
	})
	if err != nil {
		log.Fatalf("initRedisConfig NacosClient.GetConfig err: %v", err)
	}

	err = yaml.Unmarshal([]byte(content), &Conf.Database.RedisConf)
	if err != nil {
		log.Fatalf("initRedisConfig yaml.Unmarshal err: %v", err)
	}
}
