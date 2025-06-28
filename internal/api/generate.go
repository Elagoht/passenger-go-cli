package api

import (
	"passenger-go-cli/internal/schemas"
	"strconv"
)

func GeneratePassphrase(length int) (string, error) {
	response, _, err := Get[schemas.GenerateNewResponse](
		"/generate/new?length=" + strconv.Itoa(length),
	)
	if err != nil {
		return "", err
	}

	return response.Generated, nil
}

func AlternatePassphrase(passphrase string) (string, error) {
	request := map[string]string{
		"passphrase": passphrase,
	}

	response, _, err := Post[schemas.GenerateAlternativeResponse](
		"/generate/alternative",
		request,
	)
	if err != nil {
		return "", err
	}

	return response.Alternative, nil
}
