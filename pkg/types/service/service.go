package service

type Service struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	URLs        []string `json:"urls" yaml:"urls"`
	PublicKey   string   `json:"public_key" yaml:"public_key"`
}
