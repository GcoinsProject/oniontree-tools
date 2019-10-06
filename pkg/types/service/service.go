package service

type Service struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URLs        []string `json:"urls"`
}
