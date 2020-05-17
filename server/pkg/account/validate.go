package account

import "github.com/liam923/Kript/server/pkg/proto/kript/api"

func (s *Server) validateAccessTokenFormat(accessToken *api.AccessToken) bool {
	return accessToken != nil && accessToken.Jwt != nil
}

func (s *Server) validateRefreshTokenFormat(refreshToken *api.RefreshToken) bool {
	return refreshToken != nil && refreshToken.Jwt != nil
}
