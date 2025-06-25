package api

import (
	"encoding/json"
	"fmt"
	"io"
	"passenger-go-cli/internal/schemas"
)

func Login(passphrase string) (string, error) {
	client := GetClient()
	resp, err := client.Post(
		"/api/auth/login",
		map[string]string{"passphrase": passphrase},
	)

	if err != nil {
		return "", err
	}

	// Read the body and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response schemas.ResponseLogin
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	if response.Token == "" {
		var errorResponse schemas.ResponseError
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return "", fmt.Errorf("%s", errorResponse.Message)
		}
		return "", fmt.Errorf("%s", errorResponse.Message)
	}
	return response.Token, nil
}
