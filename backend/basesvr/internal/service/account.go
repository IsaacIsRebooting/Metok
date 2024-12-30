package service

import (
	"context"

	pb "github.com/IsaacIsRebooting/Metok/backend/basesvr/api/basesvr/v1"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/utils"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/service/svciface"
	"github.com/go-kratos/kratos/v2/log"
)

type AccountService struct {
	pb.UnimplementedAccountServer
	// 不能用*来传递参数？
	accountUC svciface.AccountUseCase
}

func NewAccountService(accountUC svciface.AccountUseCase) *AccountService {
	return &AccountService{
		accountUC: accountUC,
	}
}

func (s *AccountService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	// 调用uc接口创建账户逻辑
	accountID, err := s.accountUC.Create(ctx, req.GetMobile(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return &pb.RegisterReply{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}
	return &pb.RegisterReply{
		Meta:      utils.GetSuccessMeta(),
		AccountId: accountID,
	}, nil
}

func (s *AccountService) checkPassword(checkFunc func() (int64, error)) (*pb.CheckReply, error) {
	accountId, err := checkFunc()
	if err != nil {
		return &pb.CheckReply{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}
	return &pb.CheckReply{
		Meta:      utils.GetSuccessMeta(),
		AccountId: accountId,
	}, nil
}

func (s *AccountService) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckReply, error) {
	if req.GetAccountId() != 0 {
		return s.checkPassword(func() (int64, error) {
			return s.accountUC.CheckPasswordById(ctx, req.GetAccountId(), req.GetPassword())
		})
	}
	if req.GetMobile() != "" {
		return s.checkPassword(func() (int64, error) {
			return s.accountUC.CheckPasswordByMobile(ctx, req.GetMobile(), req.GetPassword())
		})
	}
	if req.GetEmail() != "" {
		return s.checkPassword(func() (int64, error) {
			return s.accountUC.CheckPasswordByEmail(ctx, req.GetEmail(), req.GetPassword())
		})
	}
	log.Context(ctx).Error("unkown account type")
	return &pb.CheckReply{
		Meta: utils.GetMetaWithErrorString("unkown account type"),
	}, nil
}
func (s *AccountService) Bind(ctx context.Context, req *pb.BindRequest) (*pb.BindReply, error) {

	return &pb.BindReply{}, nil
}
func (s *AccountService) Unbind(ctx context.Context, req *pb.UnbindRequest) (*pb.UnbindReply, error) {
	return &pb.UnbindReply{}, nil
}
