package config

import "go.uber.org/zap/zapcore"

// LocalConfig 本地的 nacos 配置
type LocalConfig struct {
	Nacos   Nacos  `yaml:"nacos"`
	AppName string `yaml:"appName"`
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
	Database   Database                `yaml:"db"`
	MetricPort int                     `yaml:"metricPort"`
	HttpPort   int                     `yaml:"httpPort"`
	ZapLog     ZapLogConf              `yaml:"zapLog"`
	Etcd       EtcdConf                `yaml:"etcd"`
	Grpc       GrpcConf                `yaml:"grpc"`
	Jwt        JwtConf                 `yaml:"jwt"`
	Domain     map[string]Domain       `yaml:"domain"`
	Services   map[string]ServicesConf `yaml:"services"`
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
type EtcdConf struct {
	Addrs       []string       `yaml:"addrs"`
	RWTimeout   int            `yaml:"rwTimeout"`
	DialTimeout int            `yaml:"dialTimeout"`
	Register    RegisterServer `yaml:"register"`
}
type GrpcConf struct {
	Host string `yaml:"host"`
	Port uint64 `yaml:"port"`
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
