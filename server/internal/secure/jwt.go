package secure

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

const (
	IssuerId              = "kript.api.account"
	RefreshTokenType      = "refresh"
	RefreshTokenLife      = time.Hour * 24 * 365 * 100
	AccessTokenType       = "access"
	AccessTokenLife       = time.Hour * 24
	VerificationTokenType = "verify"
	VerificationTokenLife = time.Minute * 30
)

// A manager for signing JWTs.
type jwtSigner struct {
	privateKey *rsa.PrivateKey
	issuerId   string
}

type JwtSigner interface {
	CreateAndSign(userId string, expires time.Time, tokenType string) (token string, tokenId string, err error)
}

func NewJwtSigner(privateKey []byte, issuerId string) (JwtSigner, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("private generate cannot be nil")
	}
	privateBlock, _ := pem.Decode(privateKey)
	if privateBlock == nil {
		return nil, fmt.Errorf("invalid private generate")
	}
	key, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private generate")
	}
	return &jwtSigner{
		privateKey: key,
		issuerId:   issuerId,
	}, nil
}

// Create and sign a JWT authenticating a user.
func (s *jwtSigner) CreateAndSign(userId string, expires time.Time, tokenType string) (token string, tokenId string, err error) {
	tokenId = uuid.New().String()
	token, err = jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    s.issuerId,
		Subject:   userId,
		Audience:  tokenType,
		ExpiresAt: expires.Unix(),
		Id:        tokenId,
	}).SignedString(s.privateKey)
	return
}

type jwtValidator struct {
	publicKey *rsa.PublicKey
	issuerId  string
}

type JwtValidator interface {
	Validate(tokenString string) (user string, tokenType string, jwtId string, err error)
}

func NewJwtValidator(publicKey []byte, issuerId string) (JwtValidator, error) {
	if publicKey == nil {
		return nil, fmt.Errorf("public generate cannot be nil")
	}
	publicBlock, _ := pem.Decode(publicKey)
	if publicBlock == nil {
		return nil, fmt.Errorf("invalid private generate")
	}
	publicParse, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid public generate")
	}
	key := publicParse.(*rsa.PublicKey)
	return &jwtValidator{
		publicKey: key,
		issuerId:  issuerId,
	}, nil
}

func (v *jwtValidator) Validate(tokenString string) (user string, tokenType string, jwtId string, err error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return v.publicKey, nil
	})

	if claims.ExpiresAt < time.Now().Unix() {
		err = fmt.Errorf("expired token")
	} else if token != nil && !token.Valid {
		err = fmt.Errorf("invalid token")
	} else if claims.Issuer != v.issuerId {
		err = fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}

	if err == nil {
		user = claims.Subject
		tokenType = claims.Audience
		jwtId = claims.Id
	}

	return
}

func GenerateJwt(signer JwtSigner, userId string, tokenType string) (valid string, invalid []string) {
	valid, _, err := signer.CreateAndSign(userId, time.Now().Add(time.Hour), tokenType)
	if err != nil {
		panic(err)
	}

	badSigner, err := NewJwtSigner(GenerateKeys(4096).Private, "issuer")
	if err != nil {
		panic(err)
	}
	invalid = []string{""}

	invalidToken, _, err := badSigner.CreateAndSign(userId, time.Now().Add(time.Hour), tokenType)
	if err != nil {
		panic(err)
	}
	invalid = append(invalid, invalidToken)

	invalidToken, _, err = signer.CreateAndSign(userId, time.Now().Add(-time.Hour), tokenType)
	if err != nil {
		panic(err)
	}
	invalid = append(invalid, invalidToken)

	invalidToken, _, err = signer.CreateAndSign(userId, time.Now().Add(time.Hour), "random-type")
	if err != nil {
		panic(err)
	}
	invalid = append(invalid, invalidToken)

	invalidToken, _, err = badSigner.CreateAndSign(userId, time.Now().Add(-time.Hour), tokenType)
	if err != nil {
		panic(err)
	}
	invalid = append(invalid, invalidToken)

	return
}
