package account

import (
	"bytes"
	"context"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const hashStrength = 12

func (s *server) UpdatePassword(ctx context.Context, request *api.UpdatePasswordRequest) (*api.UpdatePasswordResponse, error) {
	if request == nil || request.OldPassword == nil || request.NewPassword == nil || request.PrivateKey == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	userId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return nil, err
	}

	fetchedUser, err := s.database.fetchUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(fetchedUser.Password.Hash, request.OldPassword.Data) != 0 {
		return nil, status.Error(codes.InvalidArgument, "incorrect old password")
	}

	err = s.database.updateUser(ctx, userId, &user{
		Password: password{
			Hash:          request.NewPassword.Data,
			Salt:          request.NewSalt,
			HashAlgorithm: request.NewPasswordHashAlgorithm,
		},
		Keys: keys{
			PrivateKey:                    request.PrivateKey.Data,
			PrivateKeyIv:                  request.PrivateKeyIv,
			PrivateKeyKeySalt:             request.PrivateKeyKeySalt,
			PrivateKeyKeyHashAlgorithm:    request.PrivateKeyKeyHashAlgorithm,
			PrivateKeyEncryptionAlgorithm: request.PrivateKeyEncryptionAlgorithm,
		},
	})
	if err != nil {
		return nil, err
	}

	updatedUser := *fetchedUser
	updatedUser.Password.Hash = request.NewPassword.Data
	updatedUser.Password.Salt = request.NewSalt
	updatedUser.Password.HashAlgorithm = request.NewPasswordHashAlgorithm
	updatedUser.Keys.PrivateKey = request.PrivateKey.Data
	updatedUser.Keys.PrivateKeyIv = request.PrivateKeyIv
	updatedUser.Keys.PrivateKeyKeySalt = request.PrivateKeyKeySalt
	updatedUser.Keys.PrivateKeyKeyHashAlgorithm = request.PrivateKeyKeyHashAlgorithm
	updatedUser.Keys.PrivateKeyEncryptionAlgorithm = request.PrivateKeyEncryptionAlgorithm
	apiUser := updatedUser.toApiUser(userId, true)
	return &api.UpdatePasswordResponse{
		User: apiUser,
	}, nil
}

func (s *server) CreateAccount(ctx context.Context, request *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {
	if request == nil || request.Password == nil || request.PrivateKey == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	available, err := s.database.isUsernameAvailable(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, status.Errorf(codes.AlreadyExists, "username %s not available", request.Username)
	}

	hashHash, err := bcrypt.GenerateFromPassword(request.Password.Data, hashStrength)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "an internal error occurred")
	}

	user := user{
		Username: request.Username,
		Password: password{
			Hash:          hashHash,
			Salt:          request.Salt,
			HashAlgorithm: request.PasswordHashAlgorithm,
		},
		Keys: keys{
			PublicKey:                     request.PublicKey,
			PrivateKey:                    request.PrivateKey.Data,
			PrivateKeyIv:                  request.PrivateKeyIv,
			PrivateKeyEncryptionAlgorithm: request.PrivateKeyEncryptionAlgorithm,
			PrivateKeyKeySalt:             request.PrivateKeyKeySalt,
			PrivateKeyKeyHashAlgorithm:    request.PrivateKeyKeyHashAlgorithm,
			DataEncryptionAlgorithm:       request.DataEncryptionAlgorithm,
		},
		TwoFactor: make(map[string]twoFactorOption, 0),
	}
	userId, err := s.database.createUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	response, err := s.grantLogin(ctx, userId, user.toApiUser(userId, true))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.CreateAccountResponse{
		Response: response,
	}, nil
}

func (s *server) AddTwoFactor(ctx context.Context, request *api.AddTwoFactorRequest) (*api.AddTwoFactorResponse, error) {
	panic("implement me")
}

func (s *server) VerifyTwoFactor(ctx context.Context, request *api.VerifyTwoFactorRequest) (*api.VerifyTwoFactorResponse, error) {
	panic("implement me")
}
