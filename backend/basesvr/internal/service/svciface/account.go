package svciface

import "context"

//go:generate mockgen -source=account.go -destination=account_mock.go -package=svciface
type AccountUseCase interface {
	Create(ctx context.Context, mobile, email, password string) (int64, error)
	CheckPasswordById(ctx context.Context, id int64, password string) (int64, error)
	CheckPasswordByMobile(ctx context.Context, mobile, password string) (int64, error)
	CheckPasswordByEmail(ctx context.Context, email, password string) (int64, error)
	ModifyPassword(ctx context.Context, id int64, oldPassword, newPassword string) error
}
