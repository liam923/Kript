package account

import (
	"context"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUser(ctx context.Context, request *api.GetUserRequest) (response *api.GetUserResponse, err error) {
	if request == nil || (request.AccessToken != nil && !s.validateAccessTokenFormat(request.AccessToken)) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var loggedInUserId *string
	if request.AccessToken != nil {
		loggedInUserIdRaw, err := s.loginUserWithAccessToken(*request.AccessToken)
		if err != nil {
			return nil, err
		}
		loggedInUserId = &loggedInUserIdRaw
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
		err = status.Error(codes.InvalidArgument, "request.UserIdentifier must be set")
	default:
		err = status.Errorf(codes.InvalidArgument, "request.UserIdentifier has unexpected type %T", x)
	}
	if err != nil {
		return nil, err
	}

	apiUser := user.toApiUser(userId, loggedInUserId != nil && userId == *loggedInUserId)
	return &api.GetUserResponse{
		User: apiUser,
	}, nil
}
