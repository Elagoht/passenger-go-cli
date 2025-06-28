package api

import "passenger-go-cli/internal/schemas"

func GetAccounts() ([]schemas.Account, error) {
	response, _, err := Get[schemas.AccountsListResponse]("/accounts")
	if err != nil {
		return nil, err
	}

	return *response, nil
}

func GetAccountPassphrase(accountID string) (string, error) {
	endpoint := "/api/accounts/" + accountID + "/passphrase"

	response, _, err := Get[schemas.AccountPassphraseResponse](endpoint)
	if err != nil {
		return "", err
	}

	return string(*response), nil
}

func CreateAccount(account schemas.Account) (*schemas.Account, error) {
	response, _, err := Post[schemas.CreateAccountResponse](
		"/api/accounts",
		account,
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
	endpoint := "/api/accounts/" + accountID

	response, _, err := Put[schemas.Account](endpoint, account)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func DeleteAccount(accountID string) error {
	endpoint := "/api/accounts/" + accountID

	_, _, err := Delete[any](endpoint)
	return err
}
