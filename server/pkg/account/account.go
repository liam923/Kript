package account

import (
	"context"
	"fmt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
)

func (s *Server) UpdatePassword(ctx context.Context, request *api.UpdatePasswordRequest) (*api.UpdatePasswordResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}

	userId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return nil, err
	}

	fetchedUser, err := s.database.fetchUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if fetchedUser.Password.Hash != request.OldPassword {
		return nil, fmt.Errorf("incorrect old password")
	}

	err = s.database.updateUser(ctx, userId, &user{
		Password: password{
			Hash:          request.NewPassword,
			Salt:          request.NewSalt,
			HashAlgorithm: request.NewPasswordHashAlgorithm,
		},
		Keys: keys{
			PrivateKey:                    request.PrivateKey,
			PrivateKeyEncryptionAlgorithm: request.PrivateKeyEncryptionAlgorithm,
		},
	})
	if err != nil {
		return nil, err
	}

	updatedUser := *fetchedUser
	updatedUser.Password.Hash = request.NewPassword
	updatedUser.Password.Salt = request.NewSalt
	updatedUser.Password.HashAlgorithm = request.NewPasswordHashAlgorithm
	updatedUser.Keys.PrivateKey = request.PrivateKey
	updatedUser.Keys.PrivateKeyEncryptionAlgorithm = request.PrivateKeyEncryptionAlgorithm
	apiUser := updatedUser.toApiUser(userId, true)
	return &api.UpdatePasswordResponse{
		User: apiUser,
	}, nil
}

func (s *Server) CreateAccount(ctx context.Context, request *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}

	available, err := s.database.isUsernameAvailable(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, fmt.Errorf("username not available")
	}

	user := user{
		Username: request.Username,
		Password: password{
			Hash:          request.Password,
			Salt:          request.Salt,
			HashAlgorithm: request.PasswordHashAlgorithm,
		},
		Keys: keys{
			PublicKey:                     request.PublicKey,
			PrivateKey:                    request.PrivateKey,
			PrivateKeyEncryptionAlgorithm: request.PrivateKeyEncryptionAlgorithm,
			DataEncryptionAlgorithm:       request.DataEncryptionAlgorithm,
		},
		TwoFactor: make([]twoFactorOption, 0),
	}
	userId, err := s.database.createUser(ctx, &user)

	response, err := s.grantLogin(ctx, userId, user.toApiUser(userId, true))
	if err != nil {
		return nil, err
	}

	return &api.CreateAccountResponse{
		Response: response,
	}, nil
}
