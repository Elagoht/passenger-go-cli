package api

import (
	"passenger-go-cli/internal/schemas"
)

func GetAccounts() ([]schemas.Account, error) {
	response, _, err := Get[[]schemas.Account]("/accounts")
	if err != nil {
		return nil, err
	}

	return *response, nil
}

func GetAccount(accountID string) (*schemas.Account, error) {
	endpoint := "/accounts/" + accountID

	rawResponse, _, err := Get[schemas.Account](endpoint)
	if err != nil {
		return nil, err
	}

	return rawResponse, nil
}

func GetAccountPassphrase(accountID string) (string, error) {
	endpoint := "/accounts/" + accountID + "/passphrase"

	response, _, err := Get[schemas.AccountPassphraseResponse](endpoint)
	if err != nil {
		return "", err
	}

	return string(*response), nil
}

func CreateAccount(account schemas.UpsertAccountRequest) (*schemas.Account, error) {
	response, _, err := Post[schemas.CreateAccountResponse](
		"/accounts",
		map[string]string{
			"platform":   account.Platform,
			"identifier": account.Identifier,
			"passphrase": account.Passphrase,
			"url":        account.URL,
			"notes":      account.Notes,
		},
	)
	if err != nil {
		return nil, err
	}

	result := schemas.Account(*response)
	return &result, nil
}

func UpdateAccount(
	accountID string,
	account schemas.UpsertAccountRequest,
) error {
	endpoint := "/accounts/" + accountID

	_, _, err := Put[schemas.Account](endpoint, map[string]string{
		"platform":   account.Platform,
		"identifier": account.Identifier,
		"passphrase": account.Passphrase,
		"url":        account.URL,
		"notes":      account.Notes,
	})
	return err
}

func DeleteAccount(accountID string) error {
	endpoint := "/accounts/" + accountID

	_, _, err := Delete[any](endpoint)
	return err
}

func UpdateAccountPassphrase(accountID string, passphrase string) error {
	endpoint := "/accounts/" + accountID + "/passphrase"

	request := map[string]string{
		"passphrase": passphrase,
	}

	_, _, err := Put[any](endpoint, request)
	return err
}
