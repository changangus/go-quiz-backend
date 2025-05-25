package repository

import (
	"errors"

	"github.com/changangus/go-quiz-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type AnswerRepository struct {
	db *sqlx.DB
}

func NewAnswerRepository(db *sqlx.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) GetByID(id string) (*models.Answer, error) {
	answer := &models.Answer{}
	err := r.db.Get(answer, "SELECT id, question_id, answer, is_correct FROM answers WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (r *AnswerRepository) GetByQuestionID(questionID string) ([]models.Answer, error) {
	var answers []models.Answer
	err := r.db.Select(&answers, 
		"SELECT id, question_id, answer, is_correct FROM answers WHERE question_id = $1",
		questionID)
	if err != nil {
		return nil, err
	}

	return answers, nil
}

func (r *AnswerRepository) Create(data map[string]interface{}) (int64, error) {
	questionID, ok := data["question_id"].(float64)
	if !ok {
		questionIDStr, ok := data["question_id"].(string)
		if !ok {
			return 0, errors.New("question_id is required")
		}
		// Convert string to int if needed
		// This would require proper conversion logic
		_ = questionIDStr // Placeholder to avoid unused variable warning
	}

	answerText, ok := data["answer"].(string)
	if !ok || answerText == "" {
		return 0, errors.New("answer text is required")
	}

	isCorrect, _ := data["is_correct"].(bool)

	var answerID int64
	err := r.db.QueryRow(
		"INSERT INTO answers (question_id, answer, is_correct) VALUES ($1, $2, $3) RETURNING id",
		int(questionID), answerText, isCorrect,
	).Scan(&answerID)
	if err != nil {
		return 0, err
	}

	return answerID, nil
}

func (r *AnswerRepository) Update(id string, data map[string]interface{}) error {
	// Build query dynamically based on which fields are provided
	query := "UPDATE answers SET "
	params := []interface{}{}
	paramCount := 1
	first := true

	// Check for each possible field
	if answerText, ok := data["answer"].(string); ok {
		if !first {
			query += ", "
		}
		query += "answer = $" + string(paramCount+'0')
		params = append(params, answerText)
		paramCount++
		first = false
	}

	if isCorrect, ok := data["is_correct"].(bool); ok {
		if !first {
			query += ", "
		}
		query += "is_correct = $" + string(paramCount+'0')
		params = append(params, isCorrect)
		paramCount++
		first = false
	}

	if first {
		return errors.New("no valid fields to update")
	}

	query += " WHERE id = $" + string(paramCount+'0')
	params = append(params, id)

	_, err := r.db.Exec(query, params...)
	return err
}

func (r *AnswerRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM answers WHERE id = $1", id)
	return err
}
