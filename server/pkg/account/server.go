package account

import (
	"context"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
)

type Server struct {

}

func (s *Server) LoginUser(context.Context, *api.LoginUserRequest) (*api.LoginUserResponse, error) {
	panic("implement me")
}

func (s *Server) SendVerification(context.Context, *api.SendVerificationRequest) (*api.SendVerificationResponse, error) {
	panic("implement me")
}

func (s *Server) VerifyUser(context.Context, *api.VerifyUserRequest) (*api.VerifyUserResponse, error) {
	panic("implement me")
}

func (s *Server) UpdatePassword(context.Context, *api.UpdatePasswordRequest) (*api.UpdatePasswordResponse, error) {
	panic("implement me")
}

func (s *Server) CreateAccount(context.Context, *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {
	panic("implement me")
}

func (s *Server) RefreshAuth(context.Context, *api.RefreshAuthRequest) (*api.RefreshAuthResponse, error) {
	panic("implement me")
}

func (s *Server) GetUser(context.Context, *api.GetUserRequest) (*api.GetUserResponse, error) {
	panic("implement me")
}

