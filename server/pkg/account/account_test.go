package account

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/liam923/Kript/server/internal/generate"
	"github.com/liam923/Kript/server/internal/jwt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"reflect"
	"testing"
)

// A gomock matcher to assert that a user object is passed correctly.
type userMatcher struct {
	user         *user
	errorMessage string
}

func (u *userMatcher) Matches(x interface{}) bool {
	switch v := x.(type) {
	case *user:
		if u.user.Username == v.Username &&
			bytes.Compare(u.user.Password.Hash, v.Password.Hash) == 0 &&
			u.user.Password.HashAlgorithm == v.Password.HashAlgorithm &&
			bytes.Compare(u.user.Password.Salt, v.Password.Salt) == 0 &&
			bytes.Compare(u.user.Keys.PublicKey, v.Keys.PublicKey) == 0 &&
			bytes.Compare(u.user.Keys.PrivateKey, v.Keys.PrivateKey) == 0 &&
			bytes.Compare(u.user.Keys.PrivateKeyIv, v.Keys.PrivateKeyIv) == 0 &&
			u.user.Keys.PrivateKeyEncryptionAlgorithm == v.Keys.PrivateKeyEncryptionAlgorithm &&
			bytes.Compare(u.user.Keys.PrivateKeyKeySalt, v.Keys.PrivateKeyKeySalt) == 0 &&
			u.user.Keys.PrivateKeyKeyHashAlgorithm == v.Keys.PrivateKeyKeyHashAlgorithm &&
			u.user.Keys.DataEncryptionAlgorithm == v.Keys.DataEncryptionAlgorithm &&
			len(u.user.TwoFactor) == len(v.TwoFactor) {
			return true
		} else {
			u.errorMessage = fmt.Sprintf("%v", u.user)
			return false
		}
	default:
		u.errorMessage = fmt.Sprintf("incorrect type: %v", reflect.TypeOf(x))
		return false
	}
}

func (u *userMatcher) String() string {
	return u.errorMessage
}

func TestUpdatePassword(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server := createServer(t, db)

	userId := "1234567890"
	validToken, invalidTokens := generate.JWT(server.signer, userId, jwt.AccessTokenType)
	tests := []struct {
		testName string
		// the request to update the password
		request *api.UpdatePasswordRequest
		// the user being updated
		user user
		// the id of the user
		userId string
		// if true, the request password is correct
		validPassword bool
		// if true, the request object is valid
		validRequestObject bool
	}{
		// validPassword password change
		{
			testName: "valid password update 1",
			request: &api.UpdatePasswordRequest{
				AccessToken:                   &api.AccessToken{Jwt: &api.JWT{Token: validToken}},
				OldPassword:                   &api.HString{Data: []byte("spaghetti")},
				NewPassword:                   &api.HString{Data: []byte("pasta")},
				NewSalt:                       []byte("pepper"),
				NewPasswordHashAlgorithm:      0,
				PrivateKey:                    &api.EBytes{Data: []byte("re-encrypted")},
				PrivateKeyIv:                  []byte("hello world"),
				PrivateKeyKeySalt:             []byte("oregano"),
				PrivateKeyKeyHashAlgorithm:    0,
				PrivateKeyEncryptionAlgorithm: 0,
			},
			user: user{
				Username: "liam923",
				Password: password{
					Hash:          []byte("spaghetti"),
					Salt:          []byte("salt"),
					HashAlgorithm: 0,
				},
				Keys: keys{
					PublicKey:                     []byte("public"),
					PrivateKey:                    []byte("encrypted"),
					PrivateKeyIv:                  []byte("not hello world"),
					PrivateKeyKeySalt:             []byte("not oregano"),
					PrivateKeyKeyHashAlgorithm:    0,
					PrivateKeyEncryptionAlgorithm: 0,
					DataEncryptionAlgorithm:       0,
				},
				TwoFactor: nil,
			},
			userId:             userId,
			validPassword:      true,
			validRequestObject: true,
		},
		{
			testName: "valid password update 2",
			request: &api.UpdatePasswordRequest{
				AccessToken:                   &api.AccessToken{Jwt: &api.JWT{Token: validToken}},
				OldPassword:                   &api.HString{Data: []byte("spaghetti")},
				NewPassword:                   &api.HString{Data: []byte("pasta")},
				NewSalt:                       []byte("pepper"),
				NewPasswordHashAlgorithm:      0,
				PrivateKey:                    &api.EBytes{Data: []byte("re-encrypted")},
				PrivateKeyIv:                  []byte("hello world"),
				PrivateKeyKeySalt:             []byte("oregano"),
				PrivateKeyKeyHashAlgorithm:    0,
				PrivateKeyEncryptionAlgorithm: 0,
			},
			user: user{
				Username: "liam923",
				Password: password{
					Hash:          []byte("spaghetti"),
					Salt:          []byte("salt"),
					HashAlgorithm: 0,
				},
				Keys: keys{
					PublicKey:                     []byte("public"),
					PrivateKey:                    []byte("encrypted"),
					PrivateKeyIv:                  []byte("not hello world"),
					PrivateKeyKeySalt:             []byte("not oregano"),
					PrivateKeyKeyHashAlgorithm:    0,
					PrivateKeyEncryptionAlgorithm: 0,
					DataEncryptionAlgorithm:       0,
				},
				TwoFactor: []twoFactorOption{
					{
						Id:          "id",
						Type:        0,
						Destination: "123-456-7890",
					},
				},
			},
			userId:             userId,
			validPassword:      true,
			validRequestObject: true,
		},
		// wrong old password
		{
			testName: "wrong old password",
			request: &api.UpdatePasswordRequest{
				AccessToken:                   &api.AccessToken{Jwt: &api.JWT{Token: validToken}},
				OldPassword:                   &api.HString{Data: []byte("rigatoni")},
				NewPassword:                   &api.HString{Data: []byte("pasta")},
				NewSalt:                       []byte("pepper"),
				NewPasswordHashAlgorithm:      0,
				PrivateKey:                    &api.EBytes{Data: []byte("re-encrypted")},
				PrivateKeyIv:                  []byte("not hello world"),
				PrivateKeyKeySalt:             []byte("not oregano"),
				PrivateKeyKeyHashAlgorithm:    0,
				PrivateKeyEncryptionAlgorithm: 0,
			},
			user: user{
				Username: "liam923",
				Password: password{
					Hash:          []byte("spaghetti"),
					Salt:          []byte("salt"),
					HashAlgorithm: 0,
				},
				Keys: keys{
					PublicKey:                     []byte("public"),
					PrivateKey:                    []byte("encrypted"),
					PrivateKeyIv:                  []byte("not hello world"),
					PrivateKeyKeySalt:             []byte("not oregano"),
					PrivateKeyKeyHashAlgorithm:    0,
					PrivateKeyEncryptionAlgorithm: 0,
					DataEncryptionAlgorithm:       0,
				},
				TwoFactor: nil,
			},
			userId:             userId,
			validPassword:      false,
			validRequestObject: true,
		},
		// invalid request
		{
			testName: "nil request",
			request:  nil,
			user: user{
				Username: "liam923",
				Password: password{
					Hash:          []byte("spaghetti"),
					Salt:          []byte("salt"),
					HashAlgorithm: 0,
				},
				Keys: keys{
					PublicKey:                     []byte("public"),
					PrivateKey:                    []byte("encrypted"),
					PrivateKeyIv:                  []byte("not hello world"),
					PrivateKeyKeySalt:             []byte("not oregano"),
					PrivateKeyKeyHashAlgorithm:    0,
					PrivateKeyEncryptionAlgorithm: 0,
					DataEncryptionAlgorithm:       0,
				},
				TwoFactor: nil,
			},
			userId:             userId,
			validPassword:      false,
			validRequestObject: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.validRequestObject {
				db.EXPECT().
					fetchUserById(context.Background(), tt.userId).
					Return(&tt.user, nil)
				if tt.validPassword {
					db.EXPECT().
						updateUser(context.Background(), tt.userId, &userMatcher{user: &user{
							Password: password{
								Hash:          tt.request.NewPassword.Data,
								Salt:          tt.request.NewSalt,
								HashAlgorithm: tt.request.NewPasswordHashAlgorithm,
							},
							Keys: keys{
								PrivateKey:                    tt.request.PrivateKey.Data,
								PrivateKeyIv:                  tt.request.PrivateKeyIv,
								PrivateKeyKeySalt:             tt.request.PrivateKeyKeySalt,
								PrivateKeyKeyHashAlgorithm:    tt.request.PrivateKeyKeyHashAlgorithm,
								PrivateKeyEncryptionAlgorithm: tt.request.PrivateKeyEncryptionAlgorithm,
							},
						}}).
						Return(nil)
				}
			}

			response, err := server.UpdatePassword(context.Background(), tt.request)
			if !tt.validRequestObject || !tt.validPassword {
				if response != nil || err == nil {
					t.Errorf("expected error for invalid request")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for valid password update: %v", err)
				}

				expectedUser := user{
					Username: tt.user.Username,
					Password: password{
						Hash:          tt.request.NewPassword.Data,
						Salt:          tt.request.NewSalt,
						HashAlgorithm: tt.request.NewPasswordHashAlgorithm,
					},
					Keys: keys{
						PublicKey:                     tt.user.Keys.PublicKey,
						PrivateKey:                    tt.request.PrivateKey.Data,
						PrivateKeyIv:                  tt.request.PrivateKeyIv,
						PrivateKeyKeySalt:             tt.request.PrivateKeyKeySalt,
						PrivateKeyKeyHashAlgorithm:    tt.request.PrivateKeyKeyHashAlgorithm,
						PrivateKeyEncryptionAlgorithm: tt.request.PrivateKeyEncryptionAlgorithm,
						DataEncryptionAlgorithm:       tt.user.Keys.DataEncryptionAlgorithm,
					},
					TwoFactor: tt.user.TwoFactor,
				}.toApiUser(tt.userId, true)
				if expectedUser.Public.Username != response.User.Public.Username ||
					expectedUser.Public.Id != response.User.Public.Id ||
					expectedUser.Public.DataEncryptionAlgorithm != response.User.Public.DataEncryptionAlgorithm ||
					bytes.Compare(expectedUser.Public.PublicKey, response.User.Public.PublicKey) != 0 ||
					expectedUser.Public.PasswordHashAlgorithm != response.User.Public.PasswordHashAlgorithm ||
					bytes.Compare(expectedUser.Public.PasswordSalt, response.User.Public.PasswordSalt) != 0 ||
					bytes.Compare(expectedUser.Private.PrivateKey.Data, response.User.Private.PrivateKey.Data) != 0 ||
					bytes.Compare(expectedUser.Private.PrivateKeyIv, response.User.Private.PrivateKeyIv) != 0 ||
					bytes.Compare(expectedUser.Private.PrivateKeyKeySalt, response.User.Private.PrivateKeyKeySalt) != 0 ||
					expectedUser.Private.PrivateKeyKeyHashAlgorithm != response.User.Private.PrivateKeyKeyHashAlgorithm ||
					expectedUser.Private.PrivateKeyEncryptionAlgorithm != response.User.Private.PrivateKeyEncryptionAlgorithm {
					t.Errorf("invalid user returned: %v, expected: %v", response.User, expectedUser)
				}

				for _, invalidToken := range invalidTokens {
					response, err = server.UpdatePassword(context.Background(), &api.UpdatePasswordRequest{
						AccessToken:                   &api.AccessToken{Jwt: &api.JWT{Token: invalidToken}},
						OldPassword:                   tt.request.OldPassword,
						NewPassword:                   tt.request.NewPassword,
						NewSalt:                       tt.request.NewSalt,
						NewPasswordHashAlgorithm:      tt.request.NewPasswordHashAlgorithm,
						PrivateKey:                    tt.request.PrivateKey,
						PrivateKeyIv:                  tt.request.PrivateKeyIv,
						PrivateKeyKeySalt:             tt.request.PrivateKeyKeySalt,
						PrivateKeyKeyHashAlgorithm:    tt.request.PrivateKeyKeyHashAlgorithm,
						PrivateKeyEncryptionAlgorithm: tt.request.PrivateKeyEncryptionAlgorithm,
					})
					if response != nil || err == nil {
						t.Errorf("expected error for request with invalid auth")
					}
				}
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server := createServer(t, db)

	tests := []struct {
		// the request to create the account
		request api.CreateAccountRequest
		// the id of the created user
		userId string
		// if true, the requested username is taken
		usernameTaken bool
	}{
		{
			request: api.CreateAccountRequest{
				Username:                      "liam923",
				Password:                      &api.HString{Data: []byte("password")},
				Salt:                          []byte("salt"),
				PasswordHashAlgorithm:         0,
				PublicKey:                     []byte("1234567890"),
				PrivateKey:                    &api.EBytes{Data: []byte("0987654321")},
				PrivateKeyIv:                  []byte("init"),
				PrivateKeyKeySalt:             []byte("salty"),
				PrivateKeyKeyHashAlgorithm:    0,
				DataEncryptionAlgorithm:       0,
				PrivateKeyEncryptionAlgorithm: 0,
			},
			userId: "12345",
		},
		{
			request: api.CreateAccountRequest{
				Username:                      "username",
				Password:                      &api.HString{Data: []byte("asdfghjkl")},
				Salt:                          []byte("salty string"),
				PasswordHashAlgorithm:         0,
				PublicKey:                     []byte("public"),
				PrivateKey:                    &api.EBytes{Data: []byte("private")},
				PrivateKeyIv:                  []byte("init"),
				PrivateKeyKeySalt:             []byte("salty"),
				PrivateKeyKeyHashAlgorithm:    0,
				DataEncryptionAlgorithm:       0,
				PrivateKeyEncryptionAlgorithm: 0,
			},
			userId: "userid321",
		},
		{
			request: api.CreateAccountRequest{
				Username:                      "liam923",
				Password:                      &api.HString{Data: []byte("password")},
				Salt:                          []byte("salt"),
				PasswordHashAlgorithm:         0,
				PublicKey:                     []byte("1234567890"),
				PrivateKey:                    &api.EBytes{Data: []byte("0987654321")},
				PrivateKeyIv:                  []byte("init"),
				PrivateKeyKeySalt:             []byte("salty"),
				PrivateKeyKeyHashAlgorithm:    0,
				DataEncryptionAlgorithm:       0,
				PrivateKeyEncryptionAlgorithm: 0,
			},
			userId:        "12345",
			usernameTaken: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			db.EXPECT().
				isUsernameAvailable(context.Background(), tt.request.Username).
				Return(!tt.usernameTaken, nil)
			createdUser := user{
				Username: tt.request.Username,
				Password: password{
					Hash:          tt.request.Password.Data,
					Salt:          tt.request.Salt,
					HashAlgorithm: tt.request.PasswordHashAlgorithm,
				},
				Keys: keys{
					PublicKey:                     tt.request.PublicKey,
					PrivateKey:                    tt.request.PrivateKey.Data,
					PrivateKeyIv:                  tt.request.PrivateKeyIv,
					PrivateKeyEncryptionAlgorithm: tt.request.PrivateKeyEncryptionAlgorithm,
					PrivateKeyKeySalt:             tt.request.PrivateKeyKeySalt,
					PrivateKeyKeyHashAlgorithm:    tt.request.PrivateKeyKeyHashAlgorithm,
					DataEncryptionAlgorithm:       tt.request.DataEncryptionAlgorithm,
				},
				TwoFactor: nil,
			}
			if !tt.usernameTaken {
				db.EXPECT().
					createUser(context.Background(), &userMatcher{user: &createdUser}).
					Return(tt.userId, nil)
			}

			response, err := server.CreateAccount(context.Background(), &tt.request)
			if tt.usernameTaken && (response != nil || err == nil) {
				t.Errorf("should have failed to create duplicate user")
			} else if !tt.usernameTaken {
				if response == nil || err != nil {
					t.Errorf("unexpected error for validPassword creation: %v", err)
				} else {
					userId, tokenType, _, err := server.validator.ValidateJWT(response.Response.RefreshToken.Jwt.Token)
					if err != nil || userId != tt.userId || tokenType != jwt.RefreshTokenType {
						t.Errorf("Invalid refresh token: %s", response.Response.RefreshToken.Jwt.Token)
					}

					userId, tokenType, _, err = server.validator.ValidateJWT(response.Response.AccessToken.Jwt.Token)
					if err != nil || userId != tt.userId || tokenType != jwt.AccessTokenType {
						t.Errorf("Invalid access token: %s", response.Response.AccessToken.Jwt.Token)
					}

					expectedUser := createdUser.toApiUser(tt.userId, true)
					if expectedUser.Public.Username != response.Response.User.Public.Username ||
						expectedUser.Public.Id != response.Response.User.Public.Id ||
						expectedUser.Public.DataEncryptionAlgorithm != response.Response.User.Public.DataEncryptionAlgorithm ||
						bytes.Compare(expectedUser.Public.PublicKey, response.Response.User.Public.PublicKey) != 0 ||
						expectedUser.Public.PasswordHashAlgorithm != response.Response.User.Public.PasswordHashAlgorithm ||
						bytes.Compare(expectedUser.Public.PasswordSalt, response.Response.User.Public.PasswordSalt) != 0 ||
						bytes.Compare(expectedUser.Private.PrivateKey.Data, response.Response.User.Private.PrivateKey.Data) != 0 ||
						bytes.Compare(expectedUser.Private.PrivateKeyIv, response.Response.User.Private.PrivateKeyIv) != 0 ||
						bytes.Compare(expectedUser.Private.PrivateKeyKeySalt, response.Response.User.Private.PrivateKeyKeySalt) != 0 ||
						expectedUser.Private.PrivateKeyKeyHashAlgorithm != response.Response.User.Private.PrivateKeyKeyHashAlgorithm ||
						expectedUser.Private.PrivateKeyEncryptionAlgorithm != response.Response.User.Private.PrivateKeyEncryptionAlgorithm {
						t.Errorf("invalid user returned: %v, expected: %v", response.Response.User, expectedUser)
					}
				}
			}
		})
	}

	// test invalid request
	response, err := server.CreateAccount(context.Background(), nil)
	t.Run("invalid request", func(t *testing.T) {
		if response != nil || err == nil {
			t.Errorf("should have thrown error for invalid response")
		}
	})
}
