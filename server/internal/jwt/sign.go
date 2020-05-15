package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// A manager for signing JWTs.
type Signer struct {
	privateKey *rsa.PrivateKey
	issuerId   string
}

func NewSigner(privateKey []byte, issuerId string) (*Signer, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("private key cannot be nil")
	}
	privateBlock, _ := pem.Decode(privateKey)
	if privateBlock == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	key, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key")
	}
	return &Signer{
		privateKey: key,
		issuerId:   issuerId,
	}, nil
}

// Create and sign a JWT authenticating a user.
func (s *Signer) CreateAndSignJWT(userId string, expires time.Time, tokenType string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    s.issuerId,
		Subject:   userId,
		Audience:  tokenType,
		ExpiresAt: expires.Unix(),
	}).SignedString(s.privateKey)
}
