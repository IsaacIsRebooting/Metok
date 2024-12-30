package data

import (
	"context"

	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/dal/models"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/dal/query"
	"gorm.io/gen"
)

type PersistRepository struct{}

func NewPersistRepository() *PersistRepository {
	return &PersistRepository{}
}

func (r *PersistRepository) CreateAccount(ctx context.Context, account *models.Account) error {
	return query.Q.WithContext(ctx).Account.Create(account)
}

// count 在 PersistRepository 中用于获取满足特定条件的账户数量。
// 该方法主要计算未删除的账户数，可以在查询时添加额外的条件。
// 参数:
//
//	ctx - 上下文，用于取消查询。
//	conditions - 可变参数，包含查询条件。
//
// 返回值:
//
//	int64 - 满足条件的账户数量。
//	error - 查询中可能出现的错误。
func (r *PersistRepository) count(ctx context.Context, conditions ...gen.Condition) (int64, error) {
	// 添加默认条件，仅统计未删除的账户。
	conditions = append(conditions, query.Q.Account.IsDeleted.Is(false))
	// 执行查询并返回满足条件的账户数量。
	return query.Q.WithContext(ctx).Account.Where(
		conditions...,
	).Count()
}

// IsAccountExist 检查账户是否存在。
// 该方法通过执行计数查询来判断账户是否存在，而不是直接返回查询结果。
// 参数:
//
//	ctx context.Context: 上下文对象，用于传递请求范围的 deadline、取消信号、请求级值等。
//	conditions ...gen.Condition: 查询条件，使用变长参数允许灵活地指定多个条件。
//
// 返回值:
//
//	bool: 如果存在至少一个满足条件的账户，则返回true；否则返回false。
//	error: 如果查询过程中发生错误，则返回该错误。
func (r *PersistRepository) IsAccountExist(ctx context.Context, conditions ...gen.Condition) (bool, error) {
	// 执行计数查询，获取满足条件的账户数量。
	count, err := r.count(ctx, conditions...)
	if err != nil {
		// 如果查询过程中发生错误，将错误返回给调用者。
		return false, err
	}

	// 根据查询结果判断是否存在满足条件的账户。
	return count > 0, nil
}
func (r *PersistRepository) IsMobileExist(ctx context.Context, mobile string) (bool, error) {
	return r.IsAccountExist(ctx, query.Q.Account.Mobile.Eq(mobile))
}

func (r *PersistRepository) IsEmailExist(ctx context.Context, email string) (bool, error) {
	return r.IsAccountExist(ctx, query.Q.Account.Email.Eq(email))
}

func (r *PersistRepository) getFirst(ctx context.Context, conditions ...gen.Condition) (*models.Account, error) {
	conditions = append(conditions, query.Q.Account.IsDeleted.Is(false))
	account, err := query.Q.WithContext(ctx).Account.Where(
		conditions...,
	).First()
	if err != nil {
		return nil, err
	}
	return account, nil
}
func (r *PersistRepository) GetAccountById(ctx context.Context, id int64) (*models.Account, error) {
	return r.getFirst(ctx, query.Q.Account.ID.Eq(id))
}

func (r *PersistRepository) GetAccountByMobile(ctx context.Context, mobile string) (*models.Account, error) {
	return r.getFirst(ctx, query.Q.Account.Mobile.Eq(mobile))
}
func (r *PersistRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	return r.getFirst(ctx, query.Q.Account.Email.Eq(email))
}
func (r *PersistRepository) ModifyPassword(ctx context.Context, account *models.Account) error {
	_, err := query.Q.WithContext(ctx).Account.UpdateColumns(account)
	return err
}
