package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type Validator struct {
	publicKey *rsa.PublicKey
	issuerId  string
}

func NewValidator(publicKey []byte, issuerId string) (*Validator, error) {
	if publicKey == nil {
		return nil, fmt.Errorf("public key cannot be nil")
	}
	publicBlock, _ := pem.Decode(publicKey)
	if publicBlock == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	publicParse, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid public key")
	}
	key := publicParse.(*rsa.PublicKey)
	return &Validator{
		publicKey: key,
		issuerId:  issuerId,
	}, nil
}

func (v *Validator) ValidateJWT(tokenString string) (user string, tokenType string, err error) {
	claims := &jwt.StandardClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return v.publicKey, nil
	})
	if err != nil {
		return "", "", err
	} else if claims.Issuer != v.issuerId {
		return "", "", fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}

	return claims.Subject, claims.Audience, nil
}
