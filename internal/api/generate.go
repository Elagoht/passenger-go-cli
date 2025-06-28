package api

import "passenger-go-cli/internal/schemas"

// GeneratePassphrase generates a random passphrase of specified length
func GeneratePassphrase(length int) (string, error) {
	request := map[string]int{
		"length": length,
	}

	response, _, err := Post[schemas.GenerateNewResponse](
		"/api/generate",
		request,
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
		"/api/alternate",
		request,
	)
	if err != nil {
		return "", err
	}

	return response.Alternative, nil
}
