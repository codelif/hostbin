package admin

type HealthResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type DocumentResponse struct {
	Slug      string `json:"slug"`
	URL       string `json:"url"`
	SizeBytes int64  `json:"size_bytes"`
	SHA256    string `json:"sha256"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListDocumentsResponse struct {
	Documents []DocumentResponse `json:"documents"`
}

type DeleteResponse struct {
	Deleted bool   `json:"deleted"`
	Slug    string `json:"slug"`
}
