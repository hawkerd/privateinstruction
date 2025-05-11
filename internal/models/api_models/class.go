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
