package api

import "io"

func Login(passphrase string) (string, error) {
	client := GetClient()
	resp, err := client.DoWithAuth(
		"POST",
		"/api/auth/login",
		"",
		map[string]string{"passphrase": passphrase},
	)

	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
