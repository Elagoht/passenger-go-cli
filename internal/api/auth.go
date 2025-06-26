package api

import (
	"encoding/json"
	"fmt"
	"io"
	"passenger-go-cli/internal/schemas"
)

func GetStatus() (bool, error) {
	client := GetClient()
	resp, err := client.Get("/api/auth/status")
	if err != nil {
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var response schemas.ResponseStatus
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, fmt.Errorf("%s", err.Error())
	}

	return response.Status, nil
}

func Register(passphrase string) (string, error) {
	client := GetClient()
	resp, err := client.Post(
		"/api/auth/register",
		map[string]string{
			"passphrase": passphrase,
		},
	)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		var errorResponse schemas.ResponseError
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResponse.Message)
	}

	var response schemas.ResponseRegister
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %s", err.Error())
	}

	if response.Recovery == "" {
		return "", fmt.Errorf("empty recovery key received")
	}

	return response.Recovery, nil
}

func Login(passphrase string) (string, error) {
	client := GetClient()
	resp, err := client.Post(
		"/api/auth/login",
		map[string]string{"passphrase": passphrase},
	)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		var errorResponse schemas.ResponseError
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResponse.Message)
	}

	var response schemas.ResponseLogin
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %s", err.Error())
	}

	if response.Token == "" {
		return "", fmt.Errorf("empty token received")
	}

	return response.Token, nil
}
