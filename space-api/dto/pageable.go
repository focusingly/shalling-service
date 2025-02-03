package dto

type PageParam struct {
	Page *int `json:"page,omitempty"`
	Size *int `json:"size,omitempty"`
}
