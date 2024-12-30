package accountproviders

import (
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/biz"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/biz/repoiface"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/repostories/accountrepo"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/service"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/service/svciface"
	"github.com/google/wire"
)

// AccountRepoProviders 是一个用于账户仓储的 Wire 配置集合。
// 它包括了账户持久化仓储的构造函数和接口绑定。
var AccountRepoProviders = wire.NewSet(
	accountrepo.NewPersistRepository,
	wire.Bind(new(repoiface.AccountRepo), new(*accountrepo.PersistRepository)),
)

// AccountUseCaseProviders 是一个用于账户用例的 Wire 配置集合。
// 它包括了账户业务用例的构造函数和接口绑定。
var AccountUseCaseProviders = wire.NewSet(
	biz.NewAccountUseCase,
	wire.Bind(new(svciface.AccountUseCase), new(*biz.AccountUseCase)),
)

// AccountServiceProviderSet 是一个用于账户服务的 Wire 配置集合。
// 它包括了账户服务的构造函数以及账户仓储和用例的 Wire 配置集合。
var AccountServiceProviderSet = wire.NewSet(
	service.NewAccountService,
	AccountRepoProviders,
	AccountUseCaseProviders,
)
