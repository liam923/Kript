package account

import (
	"context"
	"fmt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
)

func (s *Server) GetUser(ctx context.Context, request *api.GetUserRequest) (response *api.GetUserResponse, err error) {
	if request == nil || !s.validateAccessTokenFormat(request.AccessToken) {
		return nil, fmt.Errorf("invalid request")
	}
	loggedInUserId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return
	}

	userId := ""
	var user *user = nil
	switch x := request.UserIdentifier.(type) {
	case *api.GetUserRequest_UserId:
		user, err = s.database.fetchUserById(ctx, x.UserId)
		userId = x.UserId
	case *api.GetUserRequest_Username:
		user, userId, err = s.database.fetchUserByUsername(ctx, x.Username)
	case nil:
		err = fmt.Errorf("request.UserIdentifier must be set")
	default:
		err = fmt.Errorf("request.UserIdentifier has unexpected type %T", x)
	}
	if err != nil {
		return
	}

	apiUser := user.toApiUser(userId, userId == loggedInUserId)
	return &api.GetUserResponse{
		User: apiUser,
	}, nil
}
