package dto

type Pagination struct {
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}
