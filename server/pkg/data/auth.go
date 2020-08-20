package data

import (
	"github.com/liam923/Kript/server/internal/secure"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func (s *server) validateAccessTokenFormat(accessToken *api.AccessToken) bool {
	return accessToken != nil && accessToken.Jwt != nil
}

func (s *server) validateRefreshTokenFormat(refreshToken *api.RefreshToken) bool {
	return refreshToken != nil && refreshToken.Jwt != nil
}
