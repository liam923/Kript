package data

import (
	"github.com/liam923/Kript/server/internal/jwt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) loginUserWithAccessToken(token api.AccessToken) (userId string, err error) {
	userId, tokenType, _, err := s.validator.ValidateJWT(token.Jwt.Token)
	if tokenType != jwt.AccessTokenType {
		return "", status.Errorf(codes.Unauthenticated, "incorrect token type: %s", tokenType)
	}
	if err != nil {
		return "", status.Error(codes.Unauthenticated, err.Error())
	}
	return userId, nil
}

func (s *Server) validateAccessTokenFormat(accessToken *api.AccessToken) bool {
	return accessToken != nil && accessToken.Jwt != nil
}

func (s *Server) validateRefreshTokenFormat(refreshToken *api.RefreshToken) bool {
	return refreshToken != nil && refreshToken.Jwt != nil
}
