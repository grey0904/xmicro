package center

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// LocalConfig 本地配置（用于初始化配置中心客户端）
type LocalConfig struct {
	Type    string            `yaml:"type" validate:"required,oneof=nacos etcd"` // 配置中心类型
	Nacos   NacosClientConfig `yaml:"nacos,omitempty"`                           // Nacos 客户端配置
	Etcd    EtcdClientConfig  `yaml:"etcd,omitempty"`                            // ETCD 客户端配置
	AppName string            `yaml:"appName" validate:"required"`               // 应用名称
}

// NacosClientConfig Nacos 客户端配置
type NacosClientConfig struct {
	Endpoints           []string `yaml:"endpoints"`           // Nacos 服务端地址列表
	Username            string   `yaml:"username"`            // 用户名
	Password            string   `yaml:"password"`            // 密码
	TimeoutMs           uint64   `yaml:"timeoutMs"`           // 超时时间（毫秒）
	NamespaceId         string   `yaml:"namespaceId"`         // 命名空间ID
	NotLoadCacheAtStart bool     `yaml:"notLoadCacheAtStart"` // 启动时是否加载缓存
	CacheDir            string   `yaml:"cacheDir"`            // 缓存目录
	LogDir              string   `yaml:"logDir"`              // 日志目录
	LogLevel            string   `yaml:"logLevel"`            // 日志级别
}

// EtcdClientConfig ETCD 客户端配置
type EtcdClientConfig struct {
	Endpoints   []string      `yaml:"endpoints" validate:"required,min=1"` // ETCD 节点地址
	Username    string        `yaml:"username"`                            // 用户名
	Password    string        `yaml:"password"`                            // 密码
	DialTimeout time.Duration `yaml:"dialTimeout" validate:"required"`     // 连接超时时间
}

// Config 服务相关配置，配置存放在
type Config struct {
	ServerHttp    ServerHttpConf          `yaml:"serverHttp"`
	ServerRpc     ServerRpcConf           `yaml:"serverRpc"`
	ServerMetrics ServerMetricsConf       `yaml:"serverMetrics"`
	Database      Database                `yaml:"db"`
	ZapLog        ZapLogConf              `yaml:"zapLog"`
	Nacos         NacosConf               `yaml:"nacos"`
	Etcd          EtcdConf                `yaml:"etcd"`
	Jwt           JwtConf                 `yaml:"jwt"`
	Registry      Registry                `yaml:"registry"`
	Domain        map[string]Domain       `yaml:"domain"`
	Services      map[string]ServicesConf `yaml:"services"`
}

type Registry struct {
	Kind          string         `yaml:"kind"`
	ConfigKey     string         `yaml:"configKey"`
	DeplaySeconds int64          `yaml:"deplaySeconds"`
	ServerConfig  []ServerConfig `yaml:"serverConfig"`
}

type ServerConfig struct {
	Name      string `yaml:"name"`
	ConfigKey string `yaml:"configKey"`
}

type ServerHttpConf struct {
	Host string `yaml:"host"`
	Port uint64 `yaml:"port"`
}

type ServerRpcConf struct {
	Host string `yaml:"host"`
	Port uint64 `yaml:"port"`
}

type ServerMetricsConf struct {
	Host string `yaml:"host"`
	Port uint64 `yaml:"port"`
}

type ServicesConf struct {
	Id         string `yaml:"id"`
	ClientHost string `yaml:"clientHost"`
	ClientPort int    `yaml:"clientPort"`
}
type Domain struct {
	Name        string `yaml:"name"`
	LoadBalance bool   `yaml:"loadBalance"`
}
type JwtConf struct {
	Secret string `yaml:"secret"`
	Exp    int64  `yaml:"exp"`
}
type ZapLogConf struct {
	Level zapcore.Level `yaml:"level"`
	File  string        `yaml:"file"`
}

type NacosConf struct {
	Weight float64 `yaml:"weight"`
}

type EtcdConf struct {
	Addrs       []string       `yaml:"addrs"`
	RWTimeout   int            `yaml:"rwTimeout"`
	DialTimeout int            `yaml:"dialTimeout"`
	Register    RegisterServer `yaml:"register"`
}

type RegisterServer struct {
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
