package forward

import (
	"context"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
)

type accountForwarder struct {
	client api.AccountServiceClient
}

func (af *accountForwarder) LoginUser(ctx context.Context, request *api.LoginUserRequest) (*api.LoginUserResponse, error) {
	return af.client.LoginUser(ctx, request)
}

func (af *accountForwarder) SendVerification(ctx context.Context, request *api.SendVerificationRequest) (*api.SendVerificationResponse, error) {
	return af.client.SendVerification(ctx, request)
}

func (af *accountForwarder) VerifyUser(ctx context.Context, request *api.VerifyUserRequest) (*api.VerifyUserResponse, error) {
	return af.client.VerifyUser(ctx, request)
}

func (af *accountForwarder) UpdatePassword(ctx context.Context, request *api.UpdatePasswordRequest) (*api.UpdatePasswordResponse, error) {
	return af.client.UpdatePassword(ctx, request)
}

func (af *accountForwarder) CreateAccount(ctx context.Context, request *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {
	return af.client.CreateAccount(ctx, request)
}

func (af *accountForwarder) RefreshAuth(ctx context.Context, request *api.RefreshAuthRequest) (*api.RefreshAuthResponse, error) {
	return af.client.RefreshAuth(ctx, request)
}

func (af *accountForwarder) GetUser(ctx context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
	return af.client.GetUser(ctx, request)
}

func NewAccountForwarder(client api.AccountServiceClient) api.AccountServiceServer {
	return &accountForwarder{client: client}
}
