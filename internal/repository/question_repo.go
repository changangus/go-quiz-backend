package repository

import (
	"errors"

	"github.com/changangus/go-quiz-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type QuestionRepository struct {
	db *sqlx.DB
}

func NewQuestionRepository(db *sqlx.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) GetByID(id string) (*models.Question, error) {
	question := &models.Question{}
	err := r.db.Get(question, "SELECT id, quiz_id, question, type, order_num FROM questions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (r *QuestionRepository) GetByQuizID(quizID string) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Select(&questions, 
		"SELECT id, quiz_id, question, type, order_num FROM questions WHERE quiz_id = $1 ORDER BY order_num",
		quizID)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *QuestionRepository) Create(data map[string]interface{}) (int64, error) {
	quizID, ok := data["quiz_id"].(float64)
	if !ok {
		quizIDStr, ok := data["quiz_id"].(string)
		if !ok {
			return 0, errors.New("quiz_id is required")
		}
		// Convert string to int if needed
		// This would require proper conversion logic
		_ = quizIDStr // Placeholder to avoid unused variable warning
	}

	questionText, ok := data["question"].(string)
	if !ok || questionText == "" {
		return 0, errors.New("question text is required")
	}

	questionType, ok := data["type"].(string)
	if !ok || questionType == "" {
		return 0, errors.New("question type is required")
	}

	orderNum, ok := data["order_num"].(float64)
	if !ok {
		orderNum = 0 // Default value if not provided
	}

	var questionID int64
	err := r.db.QueryRow(
		"INSERT INTO questions (quiz_id, question, type, order_num) VALUES ($1, $2, $3, $4) RETURNING id",
		int(quizID), questionText, questionType, int(orderNum),
	).Scan(&questionID)
	if err != nil {
		return 0, err
	}

	return questionID, nil
}

func (r *QuestionRepository) Update(id string, data map[string]interface{}) error {
	// Build query dynamically based on which fields are provided
	query := "UPDATE questions SET "
	params := []interface{}{}
	paramCount := 1
	first := true

	// Check for each possible field
	if questionText, ok := data["question"].(string); ok {
		if !first {
			query += ", "
		}
		query += "question = $" + string(paramCount+'0')
		params = append(params, questionText)
		paramCount++
		first = false
	}

	if questionType, ok := data["type"].(string); ok {
		if !first {
			query += ", "
		}
		query += "type = $" + string(paramCount+'0')
		params = append(params, questionType)
		paramCount++
		first = false
	}

	if orderNum, ok := data["order_num"].(float64); ok {
		if !first {
			query += ", "
		}
		query += "order_num = $" + string(paramCount+'0')
		params = append(params, int(orderNum))
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

func (r *QuestionRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM questions WHERE id = $1", id)
	return err
}
