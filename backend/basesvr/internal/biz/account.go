package biz

import (
	"context"
	"errors"

	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/biz/entity/account"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/biz/repoiface"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/dal/models"
)

var ErrAccountExist = errors.New("account exist")
var ErrAccountNotFound = errors.New("account not found")

// service
type AccountUseCase struct {
	accountRepo repoiface.AccountRepo
}

func NewAccountUseCase(accountRepo repoiface.AccountRepo) *AccountUseCase {
	return &AccountUseCase{accountRepo: accountRepo}
}

func (uc *AccountUseCase) checkAccountUnique(ctx context.Context, account *account.Account) error {
	// 如果mobile 不为空，则检查mobile是否已存在
	if account.Mobile != "" {
		exist, err := uc.accountRepo.IsMobileExist(ctx, account.Mobile)
		if err != nil {
			return errors.New("检查手机号是否已存在失败")
		}
		if exist {
			return errors.New("手机号已存在")
		}
	}

	// 如果email 不为空，则检查email是否已存在
	if account.Email != "" {
		exist, err := uc.accountRepo.IsEmailExist(ctx, account.Email)
		if err != nil {
			return errors.New("检查邮箱是否已存在失败")
		}
		if exist {
			return errors.New("邮箱已存在")
		}
	}
	return nil
}

func (uc *AccountUseCase) Create(ctx context.Context, mobile, email, password string) (int64, error) {
	// 配置账号基本信息
	account := account.NewAccount(
		account.WithMobile(mobile),
		account.WithEmail(email),
		account.WithPassword(password),
	)
	// 检查账号是否存在
	if err := uc.checkAccountUnique(ctx, account); err != nil {
		return 0, err
	}
	// 检查密码是否符合要求
	if !account.IsPasswordValid() {
		return 0, errors.New("密码不符合要求")
	}
	// 加密密码
	if err := account.EncryptPassword(); err != nil {
		return 0, errors.New("密码加密失败")
	}
	// 利用雪花算法生成唯一ID
	account.GenerateId()
	// 将账户信息保存到数据库
	if err := uc.accountRepo.CreateAccount(ctx, account.ToModel()); err != nil {
		return 0, errors.New("创建账户失败")
	}
	return account.ID, nil

}

func (uc *AccountUseCase) checkPassword(_ context.Context, getDataFunc func() (*models.Account, error), password string) (*account.Account, error) {
	// 从数据库获取到账户信息
	accountDo, err := getDataFunc()
	if err != nil {
		return nil, err
	}
	// 将model转换成account
	account := account.NewAccountWithModel(accountDo)
	// 将req中的密码和数据库中的密码进行比较，校验密码
	if err := account.CheckPassword(password); err != nil {
		return nil, err
	}
	return account, nil
}
func (uc *AccountUseCase) CheckPasswordById(ctx context.Context, id int64, password string) (int64, error) {
	// 1. 从数据库通过id获取密码
	account, err := uc.checkPassword(ctx, func() (*models.Account, error) {
		return uc.accountRepo.GetAccountById(ctx, id)
	}, password)
	if err != nil {
		return 0, err
	}
	return account.ID, nil
}
func (uc *AccountUseCase) CheckPasswordByMobile(ctx context.Context, mobile, password string) (int64, error) {
	account, err := uc.checkPassword(ctx, func() (*models.Account, error) {
		return uc.accountRepo.GetAccountByMobile(ctx, mobile)
	}, password)
	if err != nil {
		return 0, err
	}
	return account.ID, nil
}
func (uc *AccountUseCase) CheckPasswordByEmail(ctx context.Context, email, password string) (int64, error) {
	account, err := uc.checkPassword(ctx, func() (*models.Account, error) {
		return uc.accountRepo.GetAccountByEmail(ctx, email)
	}, password)
	if err != nil {
		return 0, err
	}
	return account.ID, nil

}
func (uc *AccountUseCase) ModifyPassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	// 1. 先获取当前账户信息并且验证当前旧密码是否正确，防止出现别人乱改的情况
	account, err := uc.checkPassword(ctx, func() (*models.Account, error) {
		return uc.accountRepo.GetAccountById(ctx, id)
	}, oldPassword)
	if err != nil {
		return err
	}
	// 2. 2. 得到当前账户信息后交给账户实体的方法来修改密码
	if err := account.ModifyPassword(newPassword); err != nil {
		return err
	}
	// 3. 3. 改完之后将具有新密码的账户实体持久化到数据库里,更新数据库
	if err := uc.accountRepo.ModifyPassword(ctx, account.ToModel()); err != nil {
		return err
	}
	return nil
}
