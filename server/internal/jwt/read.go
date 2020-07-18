package jwt

import "io/ioutil"

func ReadRSAKeyFile(path string) (key []byte, err error) {
	key, err = ioutil.ReadFile(path)
	return
}
