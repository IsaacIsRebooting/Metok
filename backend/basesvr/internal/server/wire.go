//go:build wireinject
// +build wireinject

package server

import (
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/server/accountproviders"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/service"
	"github.com/google/wire"
)

// initAccountService 初始化账户服务。
// 该函数使用wire.Build来组装账户服务所需的依赖，并返回账户服务的实例。
// 当前实现中，返回值为nil，这可能是因为服务实例的创建逻辑尚未实现，
// 或者wire.Build在组装过程中直接初始化了全局实例，而无需返回局部实例。
func initAccountApplication() *service.AccountService {
	wire.Build(accountproviders.AccountServiceProviderSet)
	return nil
}
