package account

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/liam923/Kript/server/internal/generate"
	"github.com/liam923/Kript/server/internal/jwt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"testing"
)

func TestGetUser(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server := createServer(t, db)

	// Create test cases
	userId1 := "test user 1" // this will be the id of the logged in user for each test case
	user1 := &user{
		Username: "liam923",
		Password: password{
			Hash:          []byte("hash"),
			Salt:          []byte("salt"),
			HashAlgorithm: 0,
		},
		Keys: keys{
			PublicKey:                     []byte("public"),
			PrivateKey:                    []byte("private"),
			PrivateKeyEncryptionAlgorithm: 0,
			DataEncryptionAlgorithm:       0,
		},
		TwoFactor: nil,
	}
	userId2 := "test user 2" // this will be the id of the other user that is fetched
	user2 := &user{
		Username: "otherdude",
		Password: password{
			Hash:          []byte("hashed"),
			Salt:          []byte("salty"),
			HashAlgorithm: 0,
		},
		Keys: keys{
			PublicKey:                     []byte("publickey"),
			PrivateKey:                    []byte("privatekey"),
			PrivateKeyEncryptionAlgorithm: 0,
			DataEncryptionAlgorithm:       0,
		},
		TwoFactor: map[string]twoFactorOption{
			"email": {
				Type:        0,
				Destination: "email@site.com",
			},
		},
	}

	validToken, invalidTokens := generate.JWT(server.signer, userId1, jwt.AccessTokenType)
	tests := []struct {
		testName string
		// The identifier of the intended user to be fetched
		userIdentifier string
		// Whether the user identifier is a username or user id
		isUsernameType bool
		// Whether the request is a fetch of the logged in user
		isSelf bool
		// The user corresponding to the requested user, or nil if it shouldn't exist
		user *user
		// The user id of the requested user
		userId string
		// Whether or not this request is being made anonymously.
		anonymous bool
	}{
		// Successful gets of different users
		{
			testName:       "get diff id",
			userIdentifier: userId2,
			isUsernameType: false,
			isSelf:         false,
			user:           user2,
			userId:         userId2,
		},
		{
			testName:       "get diff username",
			userIdentifier: user2.Username,
			isUsernameType: true,
			isSelf:         false,
			user:           user2,
			userId:         userId2,
		},
		{
			testName:       "get diff id",
			userIdentifier: userId2,
			isUsernameType: false,
			isSelf:         false,
			user:           user2,
			userId:         userId2,
			anonymous:      true,
		},
		{
			testName:       "get diff username",
			userIdentifier: user2.Username,
			isUsernameType: true,
			isSelf:         false,
			user:           user2,
			userId:         userId2,
			anonymous:      true,
		},
		// Successful gets of the logged in user
		{
			testName:       "get same id",
			userIdentifier: userId1,
			isUsernameType: false,
			isSelf:         true,
			user:           user1,
			userId:         userId1,
		},
		{
			testName:       "get same username",
			userIdentifier: user1.Username,
			isUsernameType: true,
			isSelf:         true,
			user:           user1,
			userId:         userId1,
		},
		// Fetch a nonexistent user
		{
			testName:       "nonexistent id",
			userIdentifier: "asyudilewd",
			isUsernameType: false,
			isSelf:         false,
			user:           nil,
			userId:         "123456",
		},
		{
			testName:       "nonexistent username",
			userIdentifier: "sadjnkkjsd",
			isUsernameType: true,
			isSelf:         false,
			user:           nil,
			userId:         "sadjnkkjsd",
		},
	}

	for _, tt := range tests {
		var request *api.GetUserRequest
		var unauthRequests []*api.GetUserRequest
		t.Run(tt.testName, func(t *testing.T) {
			if tt.isUsernameType {
				call := db.EXPECT().fetchUserByUsername(context.Background(), tt.userIdentifier)
				if tt.user == nil {
					call.Return(nil, "", fmt.Errorf("nonexistent user"))
				} else {
					call.Return(tt.user, tt.userId, nil)
				}
				request = &api.GetUserRequest{
					UserIdentifier: &api.GetUserRequest_Username{Username: tt.userIdentifier},
				}
				for _, invalidToken := range invalidTokens {
					unauthRequests = append(unauthRequests, &api.GetUserRequest{
						AccessToken:    &api.AccessToken{Jwt: &api.JWT{Token: invalidToken}},
						UserIdentifier: &api.GetUserRequest_Username{Username: tt.userIdentifier},
					})
				}
			} else {
				call := db.EXPECT().fetchUserById(context.Background(), tt.userIdentifier)
				if tt.user == nil {
					call.Return(nil, fmt.Errorf("nonexistent user"))
				} else {
					call.Return(tt.user, nil)
				}
				request = &api.GetUserRequest{
					UserIdentifier: &api.GetUserRequest_UserId{UserId: tt.userIdentifier},
				}
				for _, invalidToken := range invalidTokens {
					unauthRequests = append(unauthRequests, &api.GetUserRequest{
						AccessToken:    &api.AccessToken{Jwt: &api.JWT{Token: invalidToken}},
						UserIdentifier: &api.GetUserRequest_UserId{UserId: tt.userIdentifier},
					})
				}
			}
			if !tt.anonymous {
				request.AccessToken = &api.AccessToken{Jwt: &api.JWT{Token: validToken}}
			}

			response, err := server.GetUser(context.Background(), request)
			if tt.user == nil {
				if err == nil || response != nil {
					t.Errorf("expected error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error on validPassword fetch: %v", err)
				}
				expectedUser := tt.user.toApiUser(tt.userId, tt.isSelf)
				if expectedUser.Public.Username != response.User.Public.Username ||
					expectedUser.Public.Id != response.User.Public.Id ||
					expectedUser.Public.DataEncryptionAlgorithm != response.User.Public.DataEncryptionAlgorithm ||
					bytes.Compare(expectedUser.Public.PublicKey, response.User.Public.PublicKey) != 0 ||
					expectedUser.Public.PasswordHashAlgorithm != response.User.Public.PasswordHashAlgorithm ||
					bytes.Compare(expectedUser.Public.PasswordSalt, response.User.Public.PasswordSalt) != 0 ||
					(expectedUser.Private == nil && response.User.Private != nil) ||
					(expectedUser.Private != nil && (bytes.Compare(expectedUser.Private.PrivateKey.Data, response.User.Private.PrivateKey.Data) != 0 ||
						expectedUser.Private.PrivateKeyEncryptionAlgorithm != response.User.Private.PrivateKeyEncryptionAlgorithm)) {
					t.Errorf("invalid user returned: %v, expected: %v", response.User, expectedUser)
				}
			}

			for _, unauthRequest := range unauthRequests {
				response, err = server.GetUser(context.Background(), unauthRequest)
				if err == nil || response != nil {
					t.Errorf("expected error on unauthorized request")
				}
			}
		})
	}

	// Invalid request formats
	invalidRequests := []*api.GetUserRequest{
		{
			AccessToken:    nil,
			UserIdentifier: nil,
		},
		{
			AccessToken:    &api.AccessToken{Jwt: nil},
			UserIdentifier: nil,
		},
		{
			AccessToken:    &api.AccessToken{Jwt: &api.JWT{Token: validToken}},
			UserIdentifier: nil,
		},
		{
			AccessToken:    &api.AccessToken{Jwt: nil},
			UserIdentifier: &api.GetUserRequest_Username{Username: "12345"},
		},
		nil,
	}
	for i, invalidRequest := range invalidRequests {
		t.Run(fmt.Sprintf("invalid request %d", i), func(t *testing.T) {
			response, err := server.GetUser(context.Background(), invalidRequest)
			if response != nil || err == nil {
				t.Errorf("expected error")
			}
		})
	}
}
