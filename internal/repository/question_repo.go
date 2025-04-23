package repository

import (
	"database/sql"

	"github.com/changangus/go-quiz-backend/internal/models"
)

func CreateQuestion(db *sql.DB, question *models.Question) (int64, error) {
	var questionID int64
	err := db.QueryRow(
		"INSERT INTO questions (quiz_id, question, type, order_num) VALUES ($1, $2, $3, $4) RETURNING id",
		question.QuizID, question.Question, question.Type, question.Order,
	).Scan(&questionID)
	if err != nil {
		return 0, err
	}

	return questionID, nil
}
