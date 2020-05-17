package generate

import (
	"github.com/liam923/Kript/server/internal/jwt"
	"time"
)

func JWT(signer *jwt.Signer, userId string, tokenType string) (valid string, invalid []string) {
	valid, _, err := signer.CreateAndSignJWT(userId, time.Now().Add(time.Hour), tokenType)
	if err != nil {
		panic(err)
	}

	badSigner, err := jwt.NewSigner(Keys(4096).Private, "issuer")
	if err != nil {
		panic(err)
	}
	invalid = []string{""}

	invalidToken, _, err := badSigner.CreateAndSignJWT(userId, time.Now().Add(time.Hour), tokenType)
	if err != nil {
		panic(err)
	}
	invalid = append(invalid, invalidToken)

	invalidToken, _, err = signer.CreateAndSignJWT(userId, time.Now().Add(-time.Hour), tokenType)
	if err != nil {
		panic(err)
	}
	invalid = append(invalid, invalidToken)

	invalidToken, _, err = signer.CreateAndSignJWT(userId, time.Now().Add(time.Hour), "random-type")
	if err != nil {
		panic(err)
	}
	invalid = append(invalid, invalidToken)

	invalidToken, _, err = badSigner.CreateAndSignJWT(userId, time.Now().Add(-time.Hour), tokenType)
	if err != nil {
		panic(err)
	}
	invalid = append(invalid, invalidToken)

	return
}
