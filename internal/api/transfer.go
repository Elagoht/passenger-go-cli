package api

import (
	"fmt"
	"passenger-go-cli/internal/schemas"
)

func ImportCSV(filePath string) (*schemas.ImportResponse, error) {
	response, byteeee, err := PostFile[schemas.ImportResponse]("/transfer/import", filePath, "file", nil)
	fmt.Println(response)
	fmt.Println(string(byteeee))
	return response, err
}

func ExportCSV() ([]byte, error) {
	_, rawBytes, err := Post[[]byte]("/transfer/export", nil)
	if err != nil {
		return nil, err
	}

	return rawBytes, nil
}
