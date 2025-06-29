package api

func ImportCSV(filePath string) error {
	_, _, err := PostFile[any]("/transfer/import", filePath, "file", nil)
	return err
}

func ExportCSV() ([]byte, error) {
	_, rawBytes, err := Post[[]byte]("/transfer/export", nil)
	if err != nil {
		return nil, err
	}

	return rawBytes, nil
}
