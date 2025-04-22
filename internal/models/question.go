package models

type Question struct {
	ID       int    `json:"id"`
	QuizID   int    `json:"quiz_id"`
	Question string `json:"question"`
	Type     string `json:"type"`
	Order    int    `json:"order_num"`
}
