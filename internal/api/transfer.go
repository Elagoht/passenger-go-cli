package api

// ImportCSV imports accounts from a CSV file
func ImportCSV(filePath string) error {
	_, _, err := PostFile[any]("/api/import/csv", filePath, "file", nil)
	return err
}

// ExportCSV exports accounts to CSV format and returns the raw CSV data
func ExportCSV() ([]byte, error) {
	_, rawBytes, err := Get[any]("/api/export/csv")
	if err != nil {
		return nil, err
	}

	return rawBytes, nil
}
