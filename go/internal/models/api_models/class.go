package api_models

// create class
type CreateClassRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// read class
type ReadClassResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
}

// update class
type UpdateClassRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// generate join code
type GenerateJoinCodeResponse struct {
	Code         string `json:"code"`
	ExpirationDT string `json:"expiration_dt"`
}

// join class
type JoinClassRequest struct {
	JoinCode string `json:"join_code"`
}
