package models

type URLFilteredResponse struct {
	Url         string `json:"url"`
	Request     string `json:"request"`
	Description string `json:"description"`
}
