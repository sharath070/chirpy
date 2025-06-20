package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)

	// docs specify that errors are never generated
	_, err := rand.Read(key)
	if err != nil { // this is just for vibes
		return "", err
	}

	token := hex.EncodeToString(key)
	return token, nil
}
