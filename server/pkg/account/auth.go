package account

import (
	"context"
	"fmt"
	"github.com/liam923/Kript/server/internal/jwt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"time"
)

func (s *Server) LoginUser(ctx context.Context, request *api.LoginUserRequest) (*api.LoginUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
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
		err = fmt.Errorf("request.UserIdentifier must be set")
	default:
		err = fmt.Errorf("request.UserIdentifier has unexpected type %T", x)
	}
	if err != nil {
		return nil, err
	}

	if user.Password.Hash != request.Password {
		return nil, fmt.Errorf("incorrect password")
	}

	if len(user.TwoFactor) != 0 {
		token, err := s.grantVerificationToken(ctx, userId)
		if err != nil {
			return nil, err
		}

		options := make([]*api.TwoFactor, len(user.TwoFactor))
		for i, twoFactor := range user.TwoFactor {
			options[i] = &api.TwoFactor{
				Id:          twoFactor.Id,
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
			return nil, err
		}
		return &api.LoginUserResponse{
			ResponseType: &api.LoginUserResponse_Response{
				Response: response,
			},
		}, nil
	}
}

func (s *Server) SendVerification(context.Context, *api.SendVerificationRequest) (*api.SendVerificationResponse, error) {
	return nil, fmt.Errorf("two factor auth is unimplemented")
}

func (s *Server) VerifyUser(context.Context, *api.VerifyUserRequest) (*api.VerifyUserResponse, error) {
	return nil, fmt.Errorf("two factor auth is unimplemented")
}

func (s *Server) RefreshAuth(ctx context.Context, request *api.RefreshAuthRequest) (*api.RefreshAuthResponse, error) {
	if request == nil || !s.validateRefreshTokenFormat(request.RefreshToken) {
		return nil, fmt.Errorf("invalid request")
	}

	userId, tokenType, _, err := s.validator.ValidateJWT(request.RefreshToken.Jwt.Token)
	if err != nil {
		return nil, err
	}
	if tokenType != jwt.RefreshTokenType {
		return nil, fmt.Errorf("incorrect token type: %s", tokenType)
	}

	accessToken, _, err := s.signer.CreateAndSignJWT(userId, time.Now().Add(s.accessTokenLife), jwt.AccessTokenType)
	if err != nil {
		return nil, err
	}

	return &api.RefreshAuthResponse{
		AccessToken: &api.AccessToken{
			Jwt: &api.JWT{Token: accessToken},
		},
	}, nil
}

func (s *Server) loginUserWithAccessToken(token api.AccessToken) (userId string, err error) {
	userId, tokenType, _, err := s.validator.ValidateJWT(token.Jwt.Token)
	if tokenType != jwt.AccessTokenType {
		err = fmt.Errorf("incorrect token type: %s", tokenType)
	}
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (s *Server) grantLogin(ctx context.Context, userId string, user *api.User) (message *api.SuccessfulLoginMessage, err error) {
	// TODO: refresh tokens should be revocable
	refreshToken, _, err := s.signer.CreateAndSignJWT(userId, time.Now().Add(s.refreshTokenLife), jwt.RefreshTokenType)
	if err != nil {
		return
	}
	accessToken, _, err := s.signer.CreateAndSignJWT(userId, time.Now().Add(s.accessTokenLife), jwt.AccessTokenType)
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

func (s *Server) grantAccessToken(ctx context.Context, userId string) (token string, err error) {
	token, _, err = s.signer.CreateAndSignJWT(userId, time.Now().Add(s.accessTokenLife), jwt.AccessTokenType)
	return
}

func (s *Server) grantVerificationToken(ctx context.Context, userId string) (token string, err error) {
	token, _, err = s.signer.CreateAndSignJWT(userId, time.Now().Add(s.verificationTokenLife), jwt.VerificationTokenType)
	return
}
