package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var key string = "0123456789abcdef0123456789abcdef"

func EncryptPass(pass string) (string, error) {

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	gsm, err := cipher.NewGCM(block)

	if err != nil {
		return "", err
	}

	nonce := make([]byte, gsm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)

	if err != nil {
		return "", err
	}

	password := gsm.Seal(nonce, nonce, []byte(pass), nil)

	return base64.StdEncoding.EncodeToString(password), nil

}

func DecryptPass(envPass string) (string, error) {

	envPassHash, err := base64.StdEncoding.DecodeString(envPass)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	gsm, err := cipher.NewGCM(block)

	if err != nil {
		return "", err
	}

	nonceSize := gsm.NonceSize()

	nonce, cipherText := envPassHash[:nonceSize], envPassHash[nonceSize:]

	byteDecryptedPass, err := gsm.Open(nil, nonce, cipherText, nil)

	if err != nil {
		return "", err
	}

	return string(byteDecryptedPass), err
}
