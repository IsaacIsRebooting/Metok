package models

import (
	"errors"

	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/constants"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/utils"
)

func (a *Account) IsPasswordValid(patterns ...string) bool {
	// check with the given pattern
	if len(patterns) > 0 {
		pattern := patterns[0]
		return utils.IsValidWithRegex(pattern, a.Password)
	}
	return utils.IsValidWithRegex(constants.AccountPasswordPattern, a.Password)
}
func (a *Account) EncryptPassword() error {
	if err := a.generateSalt(); err != nil {
		return err
	}
	a.Password = utils.GenerateMd5WithSalt(a.Password, a.Salt)
	return nil
}
func (a *Account) generateSalt() error {
	salt, err := utils.GetPasswordSalt()
	if err != nil {
		return err
	}
	a.Salt = salt
	return nil
}
func (a *Account) CheckPassword(password string) error {
	passwordMd5 := utils.GenerateMd5WithSalt(password, a.Salt)
	if passwordMd5 != a.Password {
		return errors.New("wrong password")
	}

	return nil
}
func (a *Account) GenerateId() {
	a.ID = utils.GetSnowflakeId()
}
