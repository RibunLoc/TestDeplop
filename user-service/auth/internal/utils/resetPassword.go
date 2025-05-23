package utils
import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomPassword(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:n], nil
}
