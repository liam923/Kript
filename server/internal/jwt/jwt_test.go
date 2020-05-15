package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/liam923/Kript/server/internal/jwt"
	"testing"
	"time"
)

type keyPair struct {
	public  []byte
	private []byte
}

type test struct {
	testName        string
	keys            keyPair
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
	tests := []test{
		// Successful sign and validate
		{
			testName: "valid 1",
			keys:            generateKeys(t, 1000),
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
			testName: "valid 2",
			keys:            generateKeys(t, 1024),
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
			testName: "valid 3",
			keys:            generateKeys(t, 2312),
			signerIssuer:    "reuonewrci",
			validatorIssuer: "reuonewrci",
			userID:          "thisisauserid-1234567890",
			tokenType:       "validate",
			expireTime:      time.Now().Add(time.Nanosecond * 1000000000000000),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        true,
		},
		{
			testName: "valid 4",
			keys:            generateKeys(t, 4096),
			signerIssuer:    "eiownx43nx9re0jx9",
			validatorIssuer: "eiownx43nx9re0jx9",
			userID:          "ecuinjekckjd",
			tokenType:       "aybdssahdjqd",
			expireTime:      time.Now().Add(time.Minute),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        true,
		},
		// Invalid public or private key
		{
			testName: "invalid key 1",
			keys:            keyPair{nil, nil},
			validPublicKey:  false,
			validPrivateKey: false,
		},
		{
			testName: "invalid key 2",
			keys:            keyPair{generateKeys(t, 512).public, nil},
			validPublicKey:  true,
			validPrivateKey: false,
		},
		{
			testName: "invalid key 3",
			keys:            keyPair{nil, generateKeys(t, 512).private},
			validPublicKey:  false,
			validPrivateKey: true,
		},
		{
			testName: "invalid key 4",
			keys:            keyPair{generateKeys(t, 4096).public, nil},
			validPublicKey:  true,
			validPrivateKey: false,
		},
		{
			testName: "invalid key 5",
			keys:            keyPair{nil, generateKeys(t, 4096).private},
			validPublicKey:  false,
			validPrivateKey: true,
		},
		{
			testName: "invalid key 6",
			keys:            keyPair{randomBytes(t, 1000), randomBytes(t, 1000)},
			validPublicKey:  false,
			validPrivateKey: false,
		},
		{
			testName: "invalid key 7",
			keys: keyPair{
				public:  randomBytes(t, len(generateKeys(t, 512).public)),
				private: randomBytes(t, len(generateKeys(t, 512).private)),
			},
			validPublicKey:  false,
			validPrivateKey: false,
		},
		// Invalid signing (non-corresponding keys)
		{
			testName: "invalid key pair 1",
			keys:            keyPair{generateKeys(t, 2048).public, generateKeys(t, 2048).private},
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
			testName: "invalid key pair 2",
			keys:            keyPair{generateKeys(t, 1024).public, generateKeys(t, 1024).private},
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
			testName: "invalid key pair 3",
			keys:            keyPair{generateKeys(t, 2312).public, generateKeys(t, 2312).private},
			signerIssuer:    "reuonewrci",
			validatorIssuer: "reuonewrci",
			userID:          "thisisauserid-1234567890",
			tokenType:       "validate",
			expireTime:      time.Now().Add(time.Nanosecond * 1000000000000000),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
		{
			testName: "invalid key pair 4",
			keys:            keyPair{generateKeys(t, 4096).public, generateKeys(t, 4096).private},
			signerIssuer:    "eiownx43nx9re0jx9",
			validatorIssuer: "eiownx43nx9re0jx9",
			userID:          "ecuinjekckjd",
			tokenType:       "aybdssahdjqd",
			expireTime:      time.Now().Add(time.Minute),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
		// Wrong issuer (issuers don't match)
		{
			testName: "invalid issuer",
			keys:            generateKeys(t, 4096),
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
			testName: "expired token 1",
			keys:            generateKeys(t, 4096),
			signerIssuer:    "kript.api",
			validatorIssuer: "kript.api",
			userID:          "user123",
			tokenType:       "access",
			expireTime:      time.Now().Add(-time.Nanosecond),
			validPublicKey:  true,
			validPrivateKey: true,
			validJWT:        false,
		},
		{
			testName: "expired token 2",
			keys:            generateKeys(t, 4096),
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
			signer, err := jwt.NewSigner(tt.keys.private, tt.signerIssuer)
			if tt.validPrivateKey && err != nil {
				t.Fatal("Failed to instantiate Signer")
				return
			} else if err == nil && !tt.validPrivateKey {
				t.Fatal("Signer initialization did not fail despite invalid key")
				return
			}
			validator, err := jwt.NewValidator(tt.keys.public, tt.validatorIssuer)
			if tt.validPublicKey && err != nil {
				t.Fatal("Failed to instantiate Validator")
				return
			} else if err == nil && !tt.validPublicKey {
				t.Fatal("Validator initialization did not fail despite invalid key")
				return
			}

			if tt.validPrivateKey && tt.validPublicKey {
				token, err := signer.CreateAndSignJWT(tt.userID, tt.expireTime, tt.tokenType)
				if err != nil {
					t.Fatalf("An unexpected error occurred during signing: %v", err)
				}
				userId, tokenType, err := validator.ValidateJWT(token)
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

func generateKeys(t *testing.T, bitSize int) keyPair {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		t.Errorf("Failed to generate test keys for bit size %d", bitSize)
		return keyPair{}
	}
	publicKey := privateKey.Public()

	private := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	publicMarshalled, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Errorf("Failed to generate test keys for bit size %d: %v", bitSize, err)
		return keyPair{}
	}
	public := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicMarshalled,
		},
	)

	return keyPair{
		public:  public,
		private: private,
	}
}
