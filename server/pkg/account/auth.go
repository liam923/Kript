package account

import (
	"context"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
)

func (s *Server) LoginUser(context context.Context, request *api.LoginUserRequest) (*api.LoginUserResponse, error) {
	panic("implement me")
}

func (s *Server) SendVerification(context.Context, *api.SendVerificationRequest) (*api.SendVerificationResponse, error) {
	panic("implement me")
}

func (s *Server) VerifyUser(context.Context, *api.VerifyUserRequest) (*api.VerifyUserResponse, error) {
	panic("implement me")
}
