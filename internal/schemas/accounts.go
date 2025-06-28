package schemas

type Account struct {
	ID         string `json:"id"`
	Platform   string `json:"platform"`
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
	Notes      string `json:"notes"`
	Strength   int    `json:"strength"`
}

type AccountsListResponse []Account

type IdentifiersResponse []string

type AccountPassphraseResponse string

type CreateAccountResponse Account

type CreateAccountRequest struct {
	Platform   string `json:"platform"`
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
	Notes      string `json:"notes"`
	Passphrase string `json:"passphrase"`
}
