package account

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/liam923/Kript/server/internal/secure"
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
	if request == nil || request.TwoFactor == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	userId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return nil, err
	}

	var code string
	switch request.TwoFactor.Type {
	case api.TwoFactorType_EMAIL:
		code, err = s.emailVerificationCodeSender.SendCode(request.TwoFactor.Destination)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	default:
		return nil, status.Error(codes.Unimplemented, "the given two factor type is not yet supported")
	}

	token, tokenId, err := s.grantVerificationToken(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.database.addVerificationTokenCode(ctx, userId, tokenId, code, &twoFactorOption{
		Type:        request.TwoFactor.Type,
		Destination: request.TwoFactor.Destination,
	})
	if err != nil {
		return nil, err
	}

	return &api.AddTwoFactorResponse{
		VerificationToken: &api.VerificationToken{
			Jwt: &api.JWT{Token: token},
		},
	}, nil
}

func (s *server) VerifyTwoFactor(ctx context.Context, request *api.VerifyTwoFactorRequest) (*api.VerifyTwoFactorResponse, error) {
	if request == nil || request.VerificationToken == nil || request.VerificationToken.Jwt == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	userId, tokenType, tokenId, err := s.validator.Validate(request.VerificationToken.Jwt.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	if tokenType != secure.VerificationTokenType {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect token type: %s", tokenType)
	}

	option, err := s.database.verifyVerificationTokenCode(ctx, userId, tokenId, request.Code)
	if err != nil {
		return nil, err
	}
	if option == nil {
		return nil, status.Error(codes.Unauthenticated, "verification token is meant for login flow")
	}

	user, err := s.database.fetchUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	user.TwoFactor[uuid.New().String()] = *option
	err = s.database.updateUser(ctx, userId, user)
	if err != nil {
		return nil, err
	}

	return &api.VerifyTwoFactorResponse{
		TwoFactor: &api.TwoFactor{
			Type:        option.Type,
			Destination: option.Destination,
		},
	}, nil
}
