package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var (
	LocalConf *LocalConfig
	Conf      *Config
	validate  = validator.New()
)

// InitConfig 加载配置
func InitConfig(appName string, configPath string) error {
	if err := loadLocalConfig(configPath); err != nil {
		return err
	}

	configCenter, err := NewConfigCenter(LocalConf.Type, getConfigCenterConfig())
	if err != nil {
		return fmt.Errorf("failed to create config center: %v", err)
	}
	defer configCenter.Close()

	return initConfigWithCenter(configCenter)
}

// loadLocalConfig 加载本地配置文件
func loadLocalConfig(configPath string) error {
	if configPath == "" {
		return fmt.Errorf("config file path is required")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", configPath)
	}

	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file failed: %v", err)
	}

	LocalConf = new(LocalConfig)
	if err := v.Unmarshal(LocalConf); err != nil {
		return fmt.Errorf("unmarshal config failed: %v", err)
	}

	return validate.Struct(LocalConf)
}

// getConfigCenterConfig 获取配置中心配置
func getConfigCenterConfig() interface{} {
	switch LocalConf.Type {
	case "nacos":
		return LocalConf.Nacos
	case "etcd":
		return LocalConf.Etcd
	default:
		return nil
	}
}

// initConfigWithCenter 从配置中心初始化配置
func initConfigWithCenter(configCenter ConfigCenter) error {
	// 获取基础配置
	content, err := configCenter.GetConfig(LocalConf.AppName + ".yaml")
	if err != nil {
		return fmt.Errorf("failed to get base config: %v", err)
	}

	if err := unmarshalConfig(content, &Conf); err != nil {
		return err
	}

	// 获取其他配置
	configs := map[string]interface{}{
		"mysql.yaml": &Conf.Database.MysqlConf,
		"redis.yaml": &Conf.Database.RedisConf,
	}

	for key, target := range configs {
		content, err = configCenter.GetConfig(key)
		if err != nil {
			return fmt.Errorf("failed to get %s config: %v", key, err)
		}
		if err = unmarshalConfig(content, target); err != nil {
			return err
		}
	}

	return ValidateConfig(Conf)
}

// unmarshalConfig 解析 YAML 配置
func unmarshalConfig(content string, target interface{}) error {
	if err := yaml.Unmarshal([]byte(content), target); err != nil {
		return fmt.Errorf("unmarshal config failed: %v", err)
	}
	return nil
}

// ValidateConfig 验证配置
func ValidateConfig(config interface{}) error {
	return validate.Struct(config)
}
