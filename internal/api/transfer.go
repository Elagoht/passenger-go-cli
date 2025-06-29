package api

import (
	"passenger-go-cli/internal/schemas"
)

func ImportCSV(filePath string) (*schemas.ImportResponse, error) {
	response, _, err := PostFile[schemas.ImportResponse]("/transfer/import", filePath, "file", nil)
	return response, err
}

func ExportCSV() ([]byte, error) {
	_, rawBytes, err := Post[[]byte]("/transfer/export", nil)
	return rawBytes, err
}
