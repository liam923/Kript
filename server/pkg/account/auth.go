package account

import (
	"context"
	"github.com/liam923/Kript/server/internal/secure"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"golang.org/x/crypto/bcrypt"

	//"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (s *server) LoginUser(ctx context.Context, request *api.LoginUserRequest) (*api.LoginUserResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userId := ""
	var user *user = nil
	var err error = nil
	switch x := request.UserIdentifier.(type) {
	case *api.LoginUserRequest_UserId:
		user, err = s.database.fetchUserById(ctx, x.UserId)
		userId = x.UserId
	case *api.LoginUserRequest_Username:
		user, userId, err = s.database.fetchUserByUsername(ctx, x.Username)
	case nil:
		err = status.Error(codes.InvalidArgument, "request.UserIdentifier must be set")
	default:
		err = status.Errorf(codes.InvalidArgument, "request.UserIdentifier has unexpected type %T", x)
	}
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword(user.Password.Hash, request.Password.Data) != nil {
		return nil, status.Error(codes.Unauthenticated, "incorrect password")
	}

	if len(user.TwoFactor) != 0 {
		token, err := s.grantVerificationToken(ctx, userId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		options := make(map[string]*api.TwoFactor, len(user.TwoFactor))
		for id, twoFactor := range user.TwoFactor {
			options[id] = &api.TwoFactor{
				Type:        twoFactor.Type,
				Destination: twoFactor.Destination,
			}
		}

		return &api.LoginUserResponse{
			ResponseType: &api.LoginUserResponse_TwoFactor{
				TwoFactor: &api.LoginUserResponse_TwoFactorInfo{
					VerificationToken: &api.VerificationToken{
						Jwt: &api.JWT{Token: token},
					},
					Options: options,
				},
			},
		}, nil
	} else {
		response, err := s.grantLogin(ctx, userId, user.toApiUser(userId, true))
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &api.LoginUserResponse{
			ResponseType: &api.LoginUserResponse_Response{
				Response: response,
			},
		}, nil
	}
}

func (s *server) SendVerification(context.Context, *api.SendVerificationRequest) (*api.SendVerificationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "two factor auth is unimplemented")
}

func (s *server) VerifyUser(context.Context, *api.VerifyUserRequest) (*api.VerifyUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "two factor auth is unimplemented")
}

func (s *server) RefreshAuth(ctx context.Context, request *api.RefreshAuthRequest) (*api.RefreshAuthResponse, error) {
	if request == nil || !s.validateRefreshTokenFormat(request.RefreshToken) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	userId, tokenType, _, err := s.validator.Validate(request.RefreshToken.Jwt.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	if tokenType != secure.RefreshTokenType {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect token type: %s", tokenType)
	}

	accessToken, _, err := s.signer.CreateAndSign(userId, time.Now().Add(s.accessTokenLife), secure.AccessTokenType)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.RefreshAuthResponse{
		AccessToken: &api.AccessToken{
			Jwt: &api.JWT{Token: accessToken},
		},
	}, nil
}

func (s *server) loginUserWithAccessToken(token api.AccessToken) (userId string, err error) {
	userId, tokenType, _, err := s.validator.Validate(token.Jwt.Token)
	if tokenType != secure.AccessTokenType {
		return "", status.Errorf(codes.Unauthenticated, "incorrect token type: %s", tokenType)
	}
	if err != nil {
		return "", status.Error(codes.Unauthenticated, err.Error())
	}
	return userId, nil
}

func (s *server) grantLogin(ctx context.Context, userId string, user *api.User) (message *api.SuccessfulLoginMessage, err error) {
	// TODO: refresh tokens should be revocable
	refreshToken, _, err := s.signer.CreateAndSign(userId, time.Now().Add(s.refreshTokenLife), secure.RefreshTokenType)
	if err != nil {
		return
	}
	accessToken, _, err := s.signer.CreateAndSign(userId, time.Now().Add(s.accessTokenLife), secure.AccessTokenType)
	if err != nil {
		return
	}

	return &api.SuccessfulLoginMessage{
		RefreshToken: &api.RefreshToken{
			Jwt: &api.JWT{Token: refreshToken},
		},
		AccessToken: &api.AccessToken{
			Jwt: &api.JWT{Token: accessToken},
		},
		User: user,
	}, nil
}

func (s *server) grantAccessToken(ctx context.Context, userId string) (token string, err error) {
	token, _, err = s.signer.CreateAndSign(userId, time.Now().Add(s.accessTokenLife), secure.AccessTokenType)
	return
}

func (s *server) grantVerificationToken(ctx context.Context, userId string) (token string, err error) {
	token, _, err = s.signer.CreateAndSign(userId, time.Now().Add(s.verificationTokenLife), secure.VerificationTokenType)
	return
}
