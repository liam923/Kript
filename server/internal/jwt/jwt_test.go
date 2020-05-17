package jwt_test

import (
	"crypto/rand"
	"fmt"
	"github.com/liam923/Kript/server/internal/generate"
	"github.com/liam923/Kript/server/internal/jwt"
	"testing"
	"time"
)

type test struct {
	testName        string
	keys            generate.Pair
	signerIssuer    string
	validatorIssuer string
	userID          string
	tokenType       string
	expireTime      time.Time
	validPublicKey  bool
	validPrivateKey bool
	validJWT        bool
}

func TestSignAndValidate(t *testing.T) {
	keys := [][]generate.Pair{
		{
			generate.Keys(1000),
			generate.Keys(1000),
		},
		{
			generate.Keys(2048),
			generate.Keys(2048),
		},
		{
			generate.Keys(4096),
			generate.Keys(4096),
		},
	}
	tests := []test{
		// Successful sign and validate
		{
			testName:        "valid 1",
			keys:            keys[0][0],
			signerIssuer:    "kript.api",
			validatorIssuer: "kript.api",
			userID:          "user123",
			tokenType:       "access",
			expireTime:      time.Now().Add(time.Hour),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        true,
		},
		{
			testName:        "valid 2",
			keys:            keys[1][0],
			signerIssuer:    "issuer",
			validatorIssuer: "issuer",
			userID:          "other user",
			tokenType:       "refresh",
			expireTime:      time.Now().Add(time.Hour * 10),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        true,
		},
		{
			testName:        "valid 3",
			keys:            keys[1][0],
			signerIssuer:    "reuonewrci",
			validatorIssuer: "reuonewrci",
			userID:          "thisisauserid-1234567890",
			tokenType:       "validate",
			expireTime:      time.Now().Add(time.Nanosecond * 1000000000000000),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        true,
		},
		// Invalid public or private generate
		{
			testName:        "invalid generate 1",
			keys:            generate.Pair{nil, nil},
			validPublicKey:  false,
			validPrivateKey: false,
		},
		{
			testName:        "invalid generate 2",
			keys:            generate.Pair{keys[0][0].Public, nil},
			validPublicKey:  true,
			validPrivateKey: false,
		},
		{
			testName:        "invalid generate 3",
			keys:            generate.Pair{nil, keys[0][0].Private},
			validPublicKey:  false,
			validPrivateKey: true,
		},
		{
			testName:        "invalid generate 4",
			keys:            generate.Pair{keys[2][0].Public, nil},
			validPublicKey:  true,
			validPrivateKey: false,
		},
		{
			testName:        "invalid generate 5",
			keys:            generate.Pair{nil, keys[2][0].Private},
			validPublicKey:  false,
			validPrivateKey: true,
		},
		{
			testName:        "invalid generate 6",
			keys:            generate.Pair{randomBytes(t, 1000), randomBytes(t, 1000)},
			validPublicKey:  false,
			validPrivateKey: false,
		},
		{
			testName: "invalid generate 7",
			keys: generate.Pair{
				Public:  randomBytes(t, len(keys[0][0].Public)),
				Private: randomBytes(t, len(keys[0][0].Private)),
			},
			validPublicKey:  false,
			validPrivateKey: false,
		},
		// Invalid signing (non-corresponding keys)
		{
			testName:        "invalid generate pair 1",
			keys:            generate.Pair{keys[0][0].Public, keys[0][1].Private},
			signerIssuer:    "kript.api",
			validatorIssuer: "kript.api",
			userID:          "user123",
			tokenType:       "access",
			expireTime:      time.Now().Add(time.Hour),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
		{
			testName:        "invalid generate pair 2",
			keys:            generate.Pair{keys[1][0].Public, keys[1][1].Private},
			signerIssuer:    "issuer",
			validatorIssuer: "issuer",
			userID:          "other user",
			tokenType:       "refresh",
			expireTime:      time.Now().Add(time.Hour * 10),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
		{
			testName:        "invalid generate pair 3",
			keys:            generate.Pair{keys[2][0].Public, keys[2][1].Private},
			signerIssuer:    "reuonewrci",
			validatorIssuer: "reuonewrci",
			userID:          "thisisauserid-1234567890",
			tokenType:       "validate",
			expireTime:      time.Now().Add(time.Nanosecond * 1000000000000000),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
		// Wrong issuer (issuers don't match)
		{
			testName:        "invalid issuer",
			keys:            keys[2][0],
			signerIssuer:    "kript.api",
			validatorIssuer: "kript.ipa",
			userID:          "user123",
			tokenType:       "access",
			expireTime:      time.Now().Add(time.Hour),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
		// Expired token
		{
			testName:        "expired token 1",
			keys:            keys[2][0],
			signerIssuer:    "kript.api",
			validatorIssuer: "kript.api",
			userID:          "user123",
			tokenType:       "access",
			expireTime:      time.Now().Add(-2 * time.Second),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
		{
			testName:        "expired token 2",
			keys:            keys[2][1],
			signerIssuer:    "kript.api",
			validatorIssuer: "kript.api",
			userID:          "user123",
			tokenType:       "access",
			expireTime:      time.Now().Add(-time.Hour),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.testName), func(t *testing.T) {
			signer, err := jwt.NewSigner(tt.keys.Private, tt.signerIssuer)
			if tt.validPrivateKey && err != nil {
				t.Fatal("Failed to instantiate Signer")
				return
			} else if err == nil && !tt.validPrivateKey {
				t.Fatal("Signer initialization did not fail despite invalid generate")
				return
			}
			validator, err := jwt.NewValidator(tt.keys.Public, tt.validatorIssuer)
			if tt.validPublicKey && err != nil {
				t.Fatal("Failed to instantiate Validator")
				return
			} else if err == nil && !tt.validPublicKey {
				t.Fatal("Validator initialization did not fail despite invalid generate")
				return
			}

			if tt.validPrivateKey && tt.validPublicKey {
				token, signedTokenId, err := signer.CreateAndSignJWT(tt.userID, tt.expireTime, tt.tokenType)
				if err != nil {
					t.Fatalf("An unexpected error occurred during signing: %v", err)
				}
				userId, tokenType, validatedTokenId, err := validator.ValidateJWT(token)
				if tt.validJWT {
					if err != nil {
						t.Fatalf("Failed to validate token that should be valid: %v", err)
					}
					if userId != tt.userID {
						t.Fatalf("Validation returned incorrect user id")
					}
					if tokenType != tt.tokenType {
						t.Fatalf("Validation returned incorrect token type")
					}
					if signedTokenId != validatedTokenId {
						t.Fatalf("Token ids %s and %s did not match", signedTokenId, validatedTokenId)
					}
				} else {
					if err == nil {
						t.Fatalf("Validation should have failed for token: %v", token)
					}
				}
			}
		})
	}
}

func randomBytes(t *testing.T, n int) []byte {
	token := make([]byte, n)
	if _, err := rand.Read(token); err != nil {
		t.Errorf("Error generating random byte sequence of lenth %d", n)
	}
	return token
}
