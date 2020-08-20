package secure

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func ReadRSAKeyFile(path string) (key []byte, err error) {
	key, err = ioutil.ReadFile(path)
	return
}

type Pair struct {
	Public  []byte
	Private []byte
}

func GenerateKeys(bitSize int) Pair {
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
