package api

import "passenger-go-cli/internal/schemas"

func Login(passphrase string) (string, error) {
	loginRequest := map[string]string{
		"passphrase": passphrase,
	}

	response, _, err := Post[schemas.ResponseLogin](
		"/auth/login",
		loginRequest,
	)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func Status() (bool, error) {
	response, _, err := Get[schemas.ResponseStatus]("/auth/status")
	if err != nil {
		return false, err
	}

	return response.Status, nil
}

func Register(passphrase string) (string, error) {
	registerRequest := map[string]string{
		"passphrase": passphrase,
	}

	response, _, err := Post[schemas.ResponseRegister](
		"/auth/register",
		registerRequest,
	)
	if err != nil {
		return "", err
	}

	return response.Recovery, nil
}

func ValidateRecovery(recoveryKey string) error {
	request := map[string]string{
		"recovery": recoveryKey,
	}

	_, _, err := Post[any]("/auth/validate", request)
	return err
}

func ChangeMasterPassphrase(passphrase string) error {
	request := map[string]string{
		"passphrase": passphrase,
	}

	_, _, err := Post[any]("/auth/passphrase", request)
	if err != nil {
		return err
	}

	return nil
}
