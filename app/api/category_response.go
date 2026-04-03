package api

// CategoryResponse is the API response shape for category payloads.
type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
