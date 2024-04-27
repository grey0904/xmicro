package nacos

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"strconv"
	"strings"
	"xmicro/internal/common/config"
)

func NewNamingClient() (naming_client.INamingClient, error) {
	var (
		sc = make([]constant.ServerConfig, 0)
		nc = config.LocalConf.Nacos
	)

	for _, value := range nc.Endpoints {
		vs := strings.Split(value, ":")

		if len(vs) < 2 {
			return nil, errors.New("endpoints configuration error")
		}

		port, err := strconv.ParseUint(vs[1], 10, 64)
		if err != nil {
			return nil, errors.New("endpoints configuration error")
		}

		sc = append(sc, constant.ServerConfig{
			IpAddr: vs[0],
			Port:   port,
		})
	}

	cc := constant.ClientConfig{
		NamespaceId:         nc.NamespaceId,
		TimeoutMs:           nc.TimeoutMs,
		NotLoadCacheAtStart: true,
		CacheDir:            nc.CacheDir,
		Username:            nc.Username,
		Password:            nc.Password,
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Nacos client: %w", err)
	}

	return client, nil
}
