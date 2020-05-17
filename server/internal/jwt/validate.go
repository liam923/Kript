package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Validator struct {
	publicKey *rsa.PublicKey
	issuerId  string
}

func NewValidator(publicKey []byte, issuerId string) (*Validator, error) {
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
	return &Validator{
		publicKey: key,
		issuerId:  issuerId,
	}, nil
}

func (v *Validator) ValidateJWT(tokenString string) (user string, tokenType string, jwtId string, err error) {
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
