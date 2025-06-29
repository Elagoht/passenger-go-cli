package schemas

type FailedOne struct {
	Platform   string `json:"platform"`
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

type ImportResponse struct {
	SuccessCount int         `json:"successCount"`
	FailedOnes   []FailedOne `json:"failedOnes"`
}

// * Note: Export endpoint returns CSV file, not JSON, so no response model needed
