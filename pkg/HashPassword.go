package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var keyHash = "%sGvgrIOVLkbxCNPeTQ8Bj0UP3x7ylAETT6CxXXZP5XRT3Ftg5zzlqZnlgx1LrrtLo\n70QVCVkvbPxxNviWN0ehRg==%s"

// keyFromAnyLength یک کلید با هر طولی رو تبدیل به ۳۲ بایت می‌کنه (برای AES-256)
func keyFromAnyLength(key string) []byte {
	ke := fmt.Sprintf(keyHash, key, key)
	hash := sha256.Sum256([]byte(ke))
	return hash[:] // 32 بایت
}

func EncryptPassword(text string, key string) (string, error) {
	block, err := aes.NewCipher(keyFromAnyLength(key))
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptPassword(cryptoText string, key string) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyFromAnyLength(key))
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
