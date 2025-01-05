package service

import (
	"context"

	pb "github.com/IsaacIsRebooting/Metok/backend/shortVideoApiService/api/svapi"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceService() *UserServiceService {
	return &UserServiceService{}
}

func (s *UserServiceService) GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest) (*pb.GetVerificationCodeResponse, error) {
	return &pb.GetVerificationCodeResponse{}, nil
}
func (s *UserServiceService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{}, nil
}
func (s *UserServiceService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{}, nil
}
func (s *UserServiceService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	return &pb.GetUserInfoResponse{}, nil
}
func (s *UserServiceService) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	return &pb.UpdateUserInfoResponse{}, nil
}
func (s *UserServiceService) BindUserVoucher(ctx context.Context, req *pb.BindUserVoucherRequest) (*pb.BindUserVoucherResponse, error) {
	return &pb.BindUserVoucherResponse{}, nil
}
func (s *UserServiceService) UnbindUserVoucher(ctx context.Context, req *pb.UnbindUserVoucherRequest) (*pb.UnbindUserVoucherResponse, error) {
	return &pb.UnbindUserVoucherResponse{}, nil
}
