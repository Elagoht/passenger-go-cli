package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "passenger-go-cli"
	tokenKey    = "jwt-token"
)

func StoreToken(token string) error {
	return keyring.Set(serviceName, tokenKey, token)
}

func GetToken() (string, error) {
	token, err := keyring.Get(serviceName, tokenKey)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve token: %w", err)
	}
	return token, nil
}

func ClearToken() error {
	return keyring.Delete(serviceName, tokenKey)
}
