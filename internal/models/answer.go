package models

type Answer struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
	Correct    bool   `json:"is_correct"`
}
