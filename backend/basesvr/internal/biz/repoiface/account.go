package repoiface

import (
	"context"

	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/dal/models"
)

//go:generate mockgen -source=account.go -destination=account_mock.go -package=repoiface
type AccountRepo interface {
	CreateAccount(ctx context.Context, account *models.Account) error
	IsMobileExist(ctx context.Context, mobile string) (bool, error)
	IsEmailExist(ctx context.Context, email string) (bool, error)
	GetAccountById(ctx context.Context, id int64) (*models.Account, error)
	GetAccountByMobile(ctx context.Context, mobile string) (*models.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	ModifyPassword(ctx context.Context, account *models.Account) error
}
