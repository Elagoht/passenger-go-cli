package schemas

type ImportResponse struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors"`
}

// * Note: Export endpoint returns CSV file, not JSON, so no response model needed
