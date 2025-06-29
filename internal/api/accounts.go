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

func CreateAccount(account schemas.CreateAccountRequest) (*schemas.Account, error) {
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
	account schemas.Account,
) (*schemas.Account, error) {
	endpoint := "/accounts/" + accountID

	response, _, err := Put[schemas.Account](endpoint, account)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func DeleteAccount(accountID string) error {
	endpoint := "/accounts/" + accountID

	_, _, err := Delete[any](endpoint)
	return err
}
