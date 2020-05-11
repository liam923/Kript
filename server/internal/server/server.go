package server

import (
	"context"
	"github.com/liam923/Kript/server/proto/kript/api"
	"sync"
)

type Backend struct {
	mu    *sync.RWMutex
}

func (b *Backend) LoginUser(context.Context, *kript_api.LoginUserRequest) (*kript_api.LoginUserResponse, error) {
	panic("implement me")
}

func (b *Backend) SendVerification(context.Context, *kript_api.SendVerificationRequest) (*kript_api.SendVerificationResponse, error) {
	panic("implement me")
}

func (b *Backend) VerifyUser(context.Context, *kript_api.VerifyUserRequest) (*kript_api.VerifyUserResponse, error) {
	panic("implement me")
}

func (b *Backend) UpdatePassword(context.Context, *kript_api.UpdatePasswordRequest) (*kript_api.UpdatePasswordResponse, error) {
	panic("implement me")
}

func (b *Backend) CreateAccount(context.Context, *kript_api.CreateAccountRequest) (*kript_api.CreateAccountResponse, error) {
	panic("implement me")
}

func (b *Backend) RefreshAuth(context.Context, *kript_api.RefreshAuthRequest) (*kript_api.RefreshAuthResponse, error) {
	panic("implement me")
}

func (b *Backend) GetUser(context.Context, *kript_api.GetUserRequest) (*kript_api.GetUserResponse, error) {
	panic("implement me")
}

func (b *Backend) GetData(context.Context, *kript_api.GetDataRequest) (*kript_api.GetDataResponse, error) {
	panic("implement me")
}

func (b *Backend) UpdateDatum(context.Context, *kript_api.UpdateDatumRequest) (*kript_api.UpdateDatumResponse, error) {
	panic("implement me")
}

func (b *Backend) CreateDatum(context.Context, *kript_api.CreateDatumRequest) (*kript_api.CreateDatumResponse, error) {
	panic("implement me")
}

func (b *Backend) DeleteDatum(context.Context, *kript_api.DeleteDatumRequest) (*kript_api.DeleteDatumResponse, error) {
	panic("implement me")
}

func (b *Backend) ShareDatum(context.Context, *kript_api.ShareDatumRequest) (*kript_api.ShareDatumResponse, error) {
	panic("implement me")
}

func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}
