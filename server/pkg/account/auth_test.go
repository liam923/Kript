package account

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/liam923/Kript/server/internal/secure"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/grpclog"
	"testing"
)

type dummyWriter struct{}

func (w *dummyWriter) Write([]byte) (n int, err error) { return }

// Create a server that writes to a dummy log.
func createServer(t *testing.T, db database) server {
	issuerId := "kript.api"
	keyPair := secure.GenerateKeys(4096)
	logger := grpclog.NewLoggerV2(&dummyWriter{}, &dummyWriter{}, &dummyWriter{})
	signer, err := secure.NewJwtSigner(keyPair.Private, issuerId)
	if err != nil {
		t.Errorf("Failed to initialize signer")
	}
	validator, err := secure.NewJwtValidator(keyPair.Public, issuerId)
	if err != nil {
		t.Errorf("Failed to initialize validator")
	}
	return server{
		database:              db,
		Logger:                &logger,
		signer:                signer,
		validator:             validator,
		refreshTokenLife:      1000,
		accessTokenLife:       2000,
		verificationTokenLife: 3000,
	}
}

func TestLoginUser(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server := createServer(t, db)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 5)
	if err != nil {
		t.Errorf("error generating password hash: %s", err)
	}

	// Create test table
	tests := []struct {
		testName string
		// The request sent in
		request api.LoginUserRequest
		// Whether or not the password is correct
		isCorrectPassword bool
		// Two factors authentication options, or nil if two factor isn't enabled
		twoFactorOptions map[string]api.TwoFactor
		// The context with which the call should be made
		ctx context.Context
		// The user that the database will return, or nil if it doesn't exist
		user *user
		// Whether or not a username or userId was given. True if username
		isUsernameType bool
		// The userId of the user
		userId string
	}{
		// Valid login
		{
			testName: "validPassword login userID",
			request: api.LoginUserRequest{
				UserIdentifier: &api.LoginUserRequest_UserId{"1234567890"},
				Password:       &api.HString{Data: []byte("password")},
			},
			isCorrectPassword: true,
			twoFactorOptions: map[string]api.TwoFactor{
				"12345": {
					Type:        0,
					Destination: "email@website.com",
				},
			},
			ctx: context.Background(),
			user: &user{
				Username: "liam923",
				Password: password{
					Hash:          hashedPassword,
					Salt:          []byte("salt"),
					HashAlgorithm: 0,
				},
				Keys: keys{
					PublicKey:                     []byte("public_key"),
					PrivateKey:                    []byte("private_key"),
					PrivateKeyEncryptionAlgorithm: 0,
					DataEncryptionAlgorithm:       0,
				},
				TwoFactor: map[string]twoFactorOption{
					"12345": {
						Type:        0,
						Destination: "email@website.com",
					},
				},
			},
			isUsernameType: false,
			userId:         "1234567890",
		},
		{
			testName: "validPassword login username",
			request: api.LoginUserRequest{
				UserIdentifier: &api.LoginUserRequest_Username{"liams923"},
				Password:       &api.HString{Data: []byte("password")},
			},
			isCorrectPassword: true,
			twoFactorOptions:  nil,
			ctx:               context.WithValue(context.Background(), "keyPair", "val"),
			user: &user{
				Username: "liams923",
				Password: password{
					Hash:          hashedPassword,
					Salt:          []byte("salty"),
					HashAlgorithm: 0,
				},
				Keys: keys{
					PublicKey:                     []byte("pbkey"),
					PrivateKey:                    []byte("prkey"),
					PrivateKeyEncryptionAlgorithm: 0,
					DataEncryptionAlgorithm:       0,
				},
				TwoFactor: nil,
			},
			isUsernameType: true,
			userId:         "userID_.z9d2/qp'as;",
		},
		// Nonexistent user
		{
			testName: "nonexistent username",
			request: api.LoginUserRequest{
				UserIdentifier: &api.LoginUserRequest_Username{"nonexistent"},
				Password:       &api.HString{Data: []byte("doesn't matter")},
			},
			ctx:            context.Background(),
			user:           nil,
			isUsernameType: true,
		},
		{
			testName: "nonexistent user id",
			request: api.LoginUserRequest{
				UserIdentifier: &api.LoginUserRequest_UserId{"wqiodnk"},
				Password:       &api.HString{Data: []byte("doesn't matter")},
			},
			ctx:            context.Background(),
			user:           nil,
			isUsernameType: false,
			userId:         "wqiodnk",
		},
		// Incorrect password
		{
			testName: "incorrect password userId",
			request: api.LoginUserRequest{
				UserIdentifier: &api.LoginUserRequest_UserId{"1234567890"},
				Password:       &api.HString{Data: []byte("passwords")},
			},
			twoFactorOptions: nil,
			ctx:              context.Background(),
			user: &user{
				Username: "liam923",
				Password: password{
					Hash:          hashedPassword,
					Salt:          []byte("salt"),
					HashAlgorithm: 0,
				},
				Keys: keys{
					PublicKey:                     []byte("public_key"),
					PrivateKey:                    []byte("private_key"),
					PrivateKeyEncryptionAlgorithm: 0,
					DataEncryptionAlgorithm:       0,
				},
				TwoFactor: nil,
			},
			isUsernameType: false,
			userId:         "1234567890",
		},
		{
			testName: "incorrect password username",
			request: api.LoginUserRequest{
				UserIdentifier: &api.LoginUserRequest_Username{"liam923"},
				Password:       &api.HString{Data: []byte("PASSWORD")},
			},
			twoFactorOptions: map[string]api.TwoFactor{
				"12345": {
					Type:        0,
					Destination: "email@website.com",
				},
			},
			ctx: context.WithValue(context.Background(), "keyPair", "val"),
			user: &user{
				Username: "liam923",
				Password: password{
					Hash:          hashedPassword,
					Salt:          []byte("salty"),
					HashAlgorithm: 0,
				},
				Keys: keys{
					PublicKey:                     []byte("pbkey"),
					PrivateKey:                    []byte("prkey"),
					PrivateKeyEncryptionAlgorithm: 0,
					DataEncryptionAlgorithm:       0,
				},
				TwoFactor: map[string]twoFactorOption{
					"12345": {
						Type:        0,
						Destination: "email@website.com",
					},
				},
			},
			isUsernameType: true,
			userId:         "userID_.z9d2/qp'as;",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.isUsernameType {
				if tt.user == nil {
					db.EXPECT().
						fetchUserByUsername(tt.ctx, tt.request.UserIdentifier.(*api.LoginUserRequest_Username).Username).
						Return(nil, "", fmt.Errorf("user not found"))
				} else {
					db.EXPECT().
						fetchUserByUsername(tt.ctx, tt.user.Username).
						Return(tt.user, tt.userId, nil)
				}
			} else {
				if tt.user == nil {
					db.EXPECT().
						fetchUserById(tt.ctx, tt.userId).
						Return(nil, fmt.Errorf("user not found"))
				} else {
					db.EXPECT().
						fetchUserById(tt.ctx, tt.userId).
						Return(tt.user, nil)
				}
			}

			response, err := server.LoginUser(tt.ctx, &tt.request)
			if tt.isCorrectPassword {
				if err != nil {
					t.Errorf("unexpected error with validPassword login: %v", err)
				}
				if tt.twoFactorOptions == nil {
					switch x := response.ResponseType.(type) {
					case *api.LoginUserResponse_Response:
						userId, tokenType, _, err := server.validator.Validate(x.Response.AccessToken.Jwt.Token)
						if err != nil || userId != tt.userId || tokenType != secure.AccessTokenType {
							t.Errorf("Invalid access token: %s", x.Response.AccessToken.Jwt.Token)
						}
						userId, tokenType, _, err = server.validator.Validate(x.Response.RefreshToken.Jwt.Token)
						if err != nil || userId != tt.userId || tokenType != secure.RefreshTokenType {
							t.Errorf("Invalid refresh token: %s", x.Response.AccessToken.Jwt.Token)
						}
					default:
						t.Errorf("response.ResponseType has unexpected type %T", x)
					}
				} else {
					switch x := response.ResponseType.(type) {
					case *api.LoginUserResponse_TwoFactor:
						userId, tokenType, _, err := server.validator.Validate(x.TwoFactor.VerificationToken.Jwt.Token)
						if err != nil || userId != tt.userId || tokenType != secure.VerificationTokenType {
							t.Errorf("Invalid access token: %s", x.TwoFactor.VerificationToken.Jwt.Token)
						}

						options := x.TwoFactor.Options
						if len(options) != len(tt.twoFactorOptions) {
							t.Errorf("Invalid two factor options: %v", options)
						}
						for id, actual := range options {
							if expected, ok := tt.twoFactorOptions[id]; !ok ||
								actual.Type != expected.Type || actual.Destination != expected.Destination {
								t.Errorf("Invalid two factor options: %v", options)
							}
						}
					default:
						t.Errorf("response.ResponseType has unexpected type %T", x)
					}
				}
			} else {
				if err == nil || response != nil {
					t.Errorf("Should have failed due to incorrect password")
				}
			}
		})

		// Test malformed requests
		malformedRequests := []*api.LoginUserRequest{
			{
				UserIdentifier: nil,
				Password:       &api.HString{Data: []byte("")},
			},
			nil,
		}
		for i, request := range malformedRequests {
			t.Run(fmt.Sprintf("malformed request%d", i), func(t *testing.T) {
				response, err := server.LoginUser(context.Background(), request)
				if response != nil || err == nil {
					t.Errorf("expected error for malformed request: %v", request)
				}
			})
		}
	}
}

func TestRefreshAuth(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server := createServer(t, db)

	userId := "user12345"
	validToken, invalidTokens := secure.GenerateJwt(server.signer, userId, secure.RefreshTokenType)

	t.Run("validPassword token", func(t *testing.T) {
		response, err := server.RefreshAuth(context.Background(), &api.RefreshAuthRequest{
			RefreshToken: &api.RefreshToken{Jwt: &api.JWT{Token: validToken}},
		})
		if err != nil {
			t.Errorf("unexpected error on validPassword token: %v", err)
		}
		actualUserId, tokenType, _, err := server.validator.Validate(response.AccessToken.Jwt.Token)
		if err != nil {
			t.Errorf("unexpected error validating validPassword token: %v", err)
		}
		if actualUserId != userId || tokenType != secure.AccessTokenType {
			t.Errorf("got incorrect access token: %v", response.AccessToken.Jwt.Token)
		}
	})

	for i, invalidToken := range invalidTokens {
		t.Run(fmt.Sprintf("invalid token%d", i), func(t *testing.T) {
			response, err := server.RefreshAuth(context.Background(), &api.RefreshAuthRequest{
				RefreshToken: &api.RefreshToken{Jwt: &api.JWT{Token: invalidToken}},
			})
			if response != nil || err == nil {
				t.Errorf("refreshing should have failed")
			}
		})
	}

	// Test malformed requests
	malformedRequests := []*api.RefreshAuthRequest{
		{
			RefreshToken: nil,
		},
		{
			RefreshToken: &api.RefreshToken{
				Jwt: nil,
			},
		},
		nil,
	}
	for i, request := range malformedRequests {
		t.Run(fmt.Sprintf("malformed request%d", i), func(t *testing.T) {
			response, err := server.RefreshAuth(context.Background(), request)
			if response != nil || err == nil {
				t.Errorf("expected error for malformed request: %v", request)
			}
		})
	}
}
