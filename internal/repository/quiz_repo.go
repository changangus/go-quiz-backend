package repository

import (
	"errors"

	"github.com/changangus/go-quiz-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type QuizRepository struct {
	db *sqlx.DB
}

func NewQuizRepository(db *sqlx.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) GetByID(id string) (*models.Quiz, error) {
	quiz := &models.Quiz{}
	err := r.db.Get(quiz, "SELECT id, title, description FROM quizzes WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return quiz, nil
}

func (r *QuizRepository) GetAll() ([]models.Quiz, error) {
	var quizzes []models.Quiz
	err := r.db.Select(&quizzes, "SELECT id, title, description FROM quizzes ORDER BY id")
	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (r *QuizRepository) Create(data map[string]interface{}) (int64, error) {
	title, ok := data["title"].(string)
	if !ok || title == "" {
		return 0, errors.New("title is required")
	}

	description, _ := data["description"].(string)

	var quizID int64
	err := r.db.QueryRow(
		"INSERT INTO quizzes (title, description) VALUES ($1, $2) RETURNING id",
		title, description,
	).Scan(&quizID)
	if err != nil {
		return 0, err
	}

	return quizID, nil
}

func (r *QuizRepository) Update(id string, data map[string]interface{}) error {
	title, titleOk := data["title"].(string)
	description, descOk := data["description"].(string)

	if !titleOk && !descOk {
		return errors.New("no valid fields to update")
	}

	// Build query dynamically based on which fields are provided
	query := "UPDATE quizzes SET "
	params := []interface{}{}
	paramCount := 1

	if titleOk {
		query += "title = $" + string(paramCount+'0')
		params = append(params, title)
		paramCount++
	}

	if descOk {
		if paramCount > 1 {
			query += ", "
		}
		query += "description = $" + string(paramCount+'0')
		params = append(params, description)
		paramCount++
	}

	query += " WHERE id = $" + string(paramCount+'0')
	params = append(params, id)

	_, err := r.db.Exec(query, params...)
	return err
}

func (r *QuizRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM quizzes WHERE id = $1", id)
	return err
}
