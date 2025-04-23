package repository

import (
	"database/sql"
	"github.com/changangus/go-quiz-backend/internal/models"
)

func GetQuizById(db *sql.DB, quizID int) (*models.Quiz, error) {
	quiz := &models.Quiz{}
	err := db.QueryRow(
		"SELECT id, title, description FROM quizzes WHERE id = $1",
		quizID,
	).Scan(&quiz.ID, &quiz.Title, &quiz.Description)
	if err != nil {
		return nil, err
	}

	return quiz, nil
}

func CreateQuiz(db *sql.DB, quiz *models.Quiz) (int64, error) {
	var quizID int64
	err := db.QueryRow(
		"INSERT INTO quizzes (title, description) VALUES ($1, $2) RETURNING id",
		quiz.Title, quiz.Description,
	).Scan(&quizID)
	if err != nil {
		return 0, err
	}

	return quizID, nil
}
