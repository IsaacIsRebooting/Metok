package test

import (
	"context"
	"testing"

	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/biz"
	repomock "github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/biz/repoiface"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/dal/models"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/utils"
	"github.com/golang/mock/gomock"
)

func init() {
	// 初始化 snowflake.Node 以避免空指针异常
	utils.InitDefaultSnowflakeNode(1)
}

func TestAccountUseCase_Create(t *testing.T) {
	// 创建一个新的控制器实例
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建一个 AccountRepo 接口的模拟实现
	mockRepo := repomock.NewMockAccountRepo(ctrl)

	// 测试用例：手机号已存在
	t.Run("MobileExists", func(t *testing.T) {
		// 设置模拟行为：当调用 IsMobileExist 时返回 true
		mockRepo.EXPECT().IsMobileExist(gomock.Any(), "1234567890").Return(true, nil).Times(1)

		// 创建 AccountUseCase 实例并传入模拟的 repo
		uc := biz.NewAccountUseCase(mockRepo)

		// 调用要测试的方法
		_, err := uc.Create(context.Background(), "1234567890", "test@example.com", "password")

		// 断言结果
		if err == nil || err.Error() != "手机号已存在" {
			t.Errorf("expected error '手机号已存在', got %v", err)
		}
	})

	// 测试用例：邮箱已存在
	t.Run("EmailExists", func(t *testing.T) {
		// 设置模拟行为：当调用 IsMobileExist 时返回 false
		mockRepo.EXPECT().IsMobileExist(gomock.Any(), "").Return(false, nil).AnyTimes()

		// 设置模拟行为：当调用 IsEmailExist 时返回 true
		mockRepo.EXPECT().IsEmailExist(gomock.Any(), "test@example.com").Return(true, nil).Times(1)

		// 创建 AccountUseCase 实例并传入模拟的 repo
		uc := biz.NewAccountUseCase(mockRepo)

		// 调用要测试的方法
		_, err := uc.Create(context.Background(), "", "test@example.com", "password")

		// 断言结果
		if err == nil || err.Error() != "邮箱已存在" {
			t.Errorf("expected error '邮箱已存在', got %v", err)
		}
	})

	// 测试用例：密码不符合要求
	t.Run("InvalidPassword", func(t *testing.T) {
		// 设置模拟行为：当调用 IsMobileExist 时返回 false
		mockRepo.EXPECT().IsMobileExist(gomock.Any(), "").Return(false, nil).AnyTimes()
		// 设置模拟行为：当调用 IsEmailExist 时返回 false
		mockRepo.EXPECT().IsEmailExist(gomock.Any(), "").Return(false, nil).AnyTimes()

		// 创建 AccountUseCase 实例并传入模拟的 repo
		uc := biz.NewAccountUseCase(mockRepo)

		// 调用要测试的方法
		_, err := uc.Create(context.Background(), "", "", "short") // 假设密码长度不足

		// 断言结果
		if err == nil || err.Error() != "密码不符合要求" {
			t.Errorf("expected error '密码不符合要求', got %v", err)
		}
	})

	// 测试用例：成功创建账户
	t.Run("Success", func(t *testing.T) {
		// 设置模拟行为：当调用 IsMobileExist 时返回 false,测试时注意anytimes以及times(1)的使用，有时候因为使用了Times(1)，导致测试失败
		// 当你使用 Times(1) 时，如果你的代码逻辑没有恰好调用 IsEmailExist 恰好一次，测试就会失败。例如：
		// 如果 IsEmailExist 根本没有被调用，那么 Times(1) 的期望就不会满足，导致测试失败。
		// 如果 IsEmailExist 被调用了多次，那么同样也会因为违反了 Times(1) 的限制而失败。
		// 而在你的具体场景中，IsEmailExist 可能根本不会被调用，因为你在测试密码不符合要求的情况下，可能在密码验证阶段就已经返回了错误，从而跳过了后续的邮箱检查。
		mockRepo.EXPECT().IsMobileExist(gomock.Any(), "1234567890").Return(false, nil).Times(1)
		// 设置模拟行为：当调用 IsEmailExist 时返回 false
		mockRepo.EXPECT().IsEmailExist(gomock.Any(), "test@example.com").Return(false, nil).AnyTimes()
		// 设置模拟行为：当调用 CreateAccount 时返回 nil 错误
		mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		// 创建 AccountUseCase 实例并传入模拟的 repo
		uc := biz.NewAccountUseCase(mockRepo)

		// 调用要测试的方法
		id, err := uc.Create(context.Background(), "1234567890", "", "valid_password")

		// 断言结果
		if err != nil {
			t.Errorf("Create failed with error: %v", err)
		}
		if id == 0 {
			t.Error("Expected a non-zero ID")
		}
	})
}

func TestAccountUseCase_CheckPasswordById(t *testing.T) {
	// 创建一个新的控制器实例
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建一个 AccountRepo 接口的模拟实现
	mockRepo := repomock.NewMockAccountRepo(ctrl)

	// 测试用例：账户不存在
	t.Run("AccountNotFound", func(t *testing.T) {
		// 设置模拟行为：当调用 GetPasswordById 时返回 nil 和错误
		mockRepo.EXPECT().GetAccountById(gomock.Any(), int64(1)).Return(nil, biz.ErrAccountNotFound).AnyTimes()

		// 创建 AccountUseCase 实例并传入模拟的 repo
		uc := biz.NewAccountUseCase(mockRepo)

		// 调用要测试的方法
		id, err := uc.CheckPasswordById(context.Background(), int64(1), "password")

		// 断言结果
		if err == nil || err != biz.ErrAccountNotFound {
			t.Errorf("expected error '%v', got %v", biz.ErrAccountNotFound, err)
		}
		if id != 0 {
			t.Errorf("expected id 0, got %v", id)
		}
	})
	// 测试用例：密码不匹配
	t.Run("PasswordMismatch", func(t *testing.T) {
		// 设置模拟行为：当调用 GetPasswordById 时返回一个账户
		mockRepo.EXPECT().GetAccountById(gomock.Any(), int64(1)).Return(&models.Account{ID: 1, Salt: "1", Password: utils.GenerateMd5WithSalt("correct_password", "1")}, nil).Times(1)

		// 创建 AccountUseCase 实例并传入模拟的 repo
		uc := biz.NewAccountUseCase(mockRepo)

		// 调用要测试的方法
		id, err := uc.CheckPasswordById(context.Background(), int64(1), "wrong_password")

		// 断言结果
		if err == nil || err.Error() != "密码不匹配" {
			t.Errorf("expected error '密码不匹配', got %v", err)
		}
		if id != 0 {
			t.Errorf("expected id 0, got %v", id)
		}
	})

	// 测试用例：成功验证密码
	t.Run("Success", func(t *testing.T) {
		// 设置模拟行为：当调用 GetPasswordById 时返回一个账户
		mockRepo.EXPECT().GetAccountById(gomock.Any(), int64(1)).Return(&models.Account{ID: int64(1), Salt: "1", Password: utils.GenerateMd5WithSalt("correct_password", "1")}, nil).AnyTimes()

		// 创建 AccountUseCase 实例并传入模拟的 repo
		uc := biz.NewAccountUseCase(mockRepo)

		// 调用要测试的方法
		id, err := uc.CheckPasswordById(context.Background(), int64(1), "correct_password")

		// 断言结果
		if err != nil {
			t.Errorf("CheckPasswordById failed with error: %v", err)
		}
		if id != 1 {
			t.Errorf("expected id 1, got %v", id)
		}
	})
}
