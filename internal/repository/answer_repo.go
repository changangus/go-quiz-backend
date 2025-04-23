package repository

import (
	"database/sql"

	"github.com/changangus/go-quiz-backend/internal/models"
)

func CreateAnswer(db *sql.DB, answer *models.Answer) error {
	_, err := db.Exec(
		"INSERT INTO answers (question_id, answer, is_correct) VALUES ($1, $2, $3)",
		answer.QuestionID, answer.Answer, answer.Correct,
	)
	return err
}
