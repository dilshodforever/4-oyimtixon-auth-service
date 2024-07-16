package service

import (
	"context"
	"log"

	pb "github.com/dilshodforever/4-oyimtixon-auth-service/genprotos/auth"
	s "github.com/dilshodforever/4-oyimtixon-auth-service/storage"
)

type AuthService struct {
	stg s.InitRoot
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(stg s.InitRoot) *AuthService {
	return &AuthService{stg: stg}
}

func (a *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res, err := a.stg.Auth().Register(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return res, nil
}

func (a *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := a.stg.Auth().Login(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return res, nil
}

func (a *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	res, err := a.stg.Auth().Logout(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return res, nil
}

func (a *AuthService) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	res, err := a.stg.Auth().ForgotPassword(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return res, nil
}

func (a *AuthService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	res, err := a.stg.Auth().ResetPassword(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return res, nil
}