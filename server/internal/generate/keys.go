package generate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type Pair struct {
	Public  []byte
	Private []byte
}

func Keys(bitSize int) Pair {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate test keys for bit size %d", bitSize))
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
		panic(fmt.Errorf("Failed to generate test keys for bit size %d: %v", bitSize, err))
	}
	public := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicMarshalled,
		},
	)

	return Pair{
		Public:  public,
		Private: private,
	}
}
