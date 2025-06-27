package api

import (
	"passenger-go-cli/internal/schemas"
)

func Status() (bool, error) {
	client, err := NewClient()
	if err != nil {
		return false, err
	}

	var response schemas.ResponseStatus
	err = client.Get("/api/auth/status", &response)
	if err != nil {
		return false, err
	}

	return response.Status, nil
}

func Register(passphrase string) (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", err
	}

	var response schemas.ResponseRegister
	err = client.Post("/api/auth/register", map[string]string{
		"passphrase": passphrase,
	}, &response)
	if err != nil {
		return "", err
	}

	return response.Recovery, nil
}

func Login(passphrase string) (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", err
	}

	var response schemas.ResponseLogin
	err = client.Post("/api/auth/login", map[string]string{
		"passphrase": passphrase,
	}, &response)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}
