package account

import (
	"errors"

	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/constants"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/dal/models"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/utils"
)

const (
	ErrInvalidPassword = "密码由大小写字母、数字、符号组成，且至少需要8位"
)

type Account struct {
	ID       int64
	Mobile   string
	Email    string
	Password string
	Salt     string
}

// NewAccount 创建并返回一个新的Account实例。
// 该函数接受可变参数options，这些参数是AccountOption类型的函数，用于配置Account实例。
// 通过这种方式，用户可以选择性地应用多个配置选项到新创建的Account实例上。
func NewAccount(options ...AccountOption) *Account {
	// 创建一个新的Account实例。
	account := &Account{}

	// 遍历所有提供的配置选项，并依次应用到account实例上。
	for _, option := range options {
		option(account)
	}

	// 返回配置完毕的Account实例。
	return account
}

func NewAccountWithModel(model *models.Account) *Account {
	return &Account{
		ID:       model.ID,
		Email:    model.Email,
		Mobile:   model.Mobile,
		Password: model.Password,
		Salt:     model.Salt,
	}
}

func (a *Account) ToModel() *models.Account {
	return &models.Account{
		ID:       a.ID,
		Email:    a.Email,
		Mobile:   a.Mobile,
		Password: a.Password,
		Salt:     a.Salt,
	}
}

func (a *Account) IsPasswordValid(pattern ...string) bool {
	if len(pattern) > 0 {
		pattern := pattern[0]
		return utils.IsValidWithRegex(pattern, a.Password)
	}
	return utils.IsValidWithRegex(constants.AccountPasswordPattern, a.Password)
}

func (a *Account) EncryptPassword() error {
	salt, err := utils.GetPasswordSalt()
	if err != nil {
		return err
	}
	a.Salt = salt
	a.Password = utils.GenerateMd5WithSalt(a.Password, salt)
	return nil
}

func (a *Account) ModifyPassword(newPassword string) error {
	// 将新密码赋值给现密码
	a.Password = newPassword
	// 判断现密码是否符合要求
	if isValid := a.IsPasswordValid(); !isValid {
		return errors.New("password is invalid")
	}
	// 加密密码
	if err := a.EncryptPassword(); err != nil {
		return err
	}
	return nil
}

func (a *Account) CheckPassword(password string) error {
	// 密码校验，先将传入的密码和盐进行加密，再和数据库中的密码进行比较
	passwordMd5 := utils.GenerateMd5WithSalt(password, a.Salt)
	if passwordMd5 != a.Password {
		return errors.New("wrong password")
	}
	return nil
}
func (a *Account) GenerateId() {
	a.ID = utils.GetSnowflakeId()
}
