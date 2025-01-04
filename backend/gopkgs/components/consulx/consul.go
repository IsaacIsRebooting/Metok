package consulx

import (
	"context"
	"fmt"
	"sync"

	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/components"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

var (
	// use sync.Map to store client, in order to deal with concurrent problem when there are multiple requests
	globalClientMap = sync.Map{}
	globalConfigMap = make(components.ConfigMap[*Config])
)

func GetConfig() components.ConfigMap[*Config] {
	return globalConfigMap
}

func Init(cm components.ConfigMap[*Config]) (func() error, error) {
	globalConfigMap = cm

	for k, v := range cm {
		client, err := Connect(v)
		if err != nil {
			return nil, err
		}
		globalClientMap.Store(k, client)
	}

	return IsHealth, nil
}

func Connect(c *Config) (*api.Client, error) {
	c.setDefault()
	cfg := api.DefaultConfig()
	cfg.Address = c.Address

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, err
}

func GetClient(ctx context.Context, keys ...string) *api.Client {
	storeKey := "default"
	if len(keys) != 0 {
		storeKey = keys[0]
	}
	client, ok := globalClientMap.Load(storeKey)
	if !ok {
		panic(fmt.Sprintf("consul client %s not found", storeKey))
	}
	return client.(*api.Client)
}

func GetGrpcConn(ctx context.Context, entryPoint string, keys ...string) (*grpc.ClientConn, error) {
	client := GetClient(ctx, keys...)
	return kratosgrpc.DialInsecure(
		ctx,
		kratosgrpc.WithEndpoint(entryPoint),
		kratosgrpc.WithDiscovery(consul.New(client)),
		kratosgrpc.WithMiddleware(
			tracing.Client(),
		),
	)
}

// IsHealth 检查所有客户端的健康状态。
// 该函数遍历全局客户端映射，对每个客户端执行健康检查。
// 如果所有客户端都通过健康检查，则返回nil，否则返回错误。
func IsHealth() (err error) {
	// 使用Range遍历globalClientMap，对每个客户端进行健康检查。
	// 注意：这里没有直接返回错误，而是将错误记录日志并继续检查其他客户端。
	globalClientMap.Range(func(key, value interface{}) bool {
		// 从映射中获取客户端。
		client := value.(*api.Client)
		// 执行健康检查，忽略检查结果，只关注是否有错误发生。
		_, _, e := client.Health().State("any", nil)
		if e != nil {
			// 如果健康检查失败，记录错误日志并停止遍历。
			log.Errorf("consul health check failed, client key: %s", key)
			return false
		}
		// 如果健康检查成功，记录信息日志并继续遍历。
		log.Infof("consul %s health check ok", key)
		return true
	})
	// 返回遍历过程中可能发生的错误。
	return err
}
