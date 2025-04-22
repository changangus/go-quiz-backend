package models

type Quiz struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
