package utils

type PageAttr struct {
	Page    int               `json:"page"`
	Size    int               `json:"size"`
	Filter  map[string]string `json:"filter"`
	SortBy  string            `json:"sort_by"`
	SortDir string            `json:"sort_dir"`
	Search  string            `json:"search"`
}
