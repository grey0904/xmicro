package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"log"
	"strconv"
	"strings"
	"xmicro/internal/app"
)

func CreateConfigClient() {
	var (
		sc = make([]constant.ServerConfig, 0)
		nc = app.Config.Nacos
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

	app.Nc = client
}
