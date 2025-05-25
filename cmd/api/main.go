package main

import (
	"log"
	"net/http"

	"github.com/changangus/go-quiz-backend/db"
	"github.com/changangus/go-quiz-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func setupRouter(database *sqlx.DB) *gin.Engine {
	router := gin.Default()

	// Create repositories
	quizRepo := repository.NewQuizRepository(database)
	questionRepo := repository.NewQuestionRepository(database)
	answerRepo := repository.NewAnswerRepository(database)

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Quizzes endpoints
		quizzes := api.Group("/quizzes")
		{
			quizzes.GET("", func(c *gin.Context) {
				quizzes, err := quizRepo.GetAll()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, quizzes)
			})

			quizzes.GET("/:id", func(c *gin.Context) {
				id := c.Param("id")
				quiz, err := quizRepo.GetByID(id)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
					return
				}
				
				// Get questions for this quiz
				questions, err := questionRepo.GetByQuizID(id)
				if err == nil {
					// If we have questions, attach them to the quiz response
					// This could be enhanced to include answers as well
					c.JSON(http.StatusOK, gin.H{
						"quiz": quiz,
						"questions": questions,
					})
					return
				}
				
				c.JSON(http.StatusOK, quiz)
			})

			quizzes.POST("", func(c *gin.Context) {
				var quiz map[string]interface{}
				if err := c.ShouldBindJSON(&quiz); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				
				id, err := quizRepo.Create(quiz)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusCreated, gin.H{"id": id})
			})
			
			quizzes.PUT("/:id", func(c *gin.Context) {
				id := c.Param("id")
				var data map[string]interface{}
				if err := c.ShouldBindJSON(&data); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				
				err := quizRepo.Update(id, data)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{"message": "Quiz updated successfully"})
			})
			
			quizzes.DELETE("/:id", func(c *gin.Context) {
				id := c.Param("id")
				err := quizRepo.Delete(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{"message": "Quiz deleted successfully"})
			})
			
			// Questions related to a quiz
			quizzes.GET("/:id/questions", func(c *gin.Context) {
				id := c.Param("id")
				questions, err := questionRepo.GetByQuizID(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, questions)
			})
			
			quizzes.POST("/:id/questions", func(c *gin.Context) {
				quizID := c.Param("id")
				var data map[string]interface{}
				if err := c.ShouldBindJSON(&data); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				
				// Make sure the quiz ID is included
				data["quiz_id"] = quizID
				
				id, err := questionRepo.Create(data)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusCreated, gin.H{"id": id})
			})
		}

		// Questions endpoints
		questions := api.Group("/questions")
		{
			questions.GET("/:id", func(c *gin.Context) {
				id := c.Param("id")
				question, err := questionRepo.GetByID(id)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
					return
				}
				
				// Get answers for this question
				answers, err := answerRepo.GetByQuestionID(id)
				if err == nil {
					c.JSON(http.StatusOK, gin.H{
						"question": question,
						"answers": answers,
					})
					return
				}
				
				c.JSON(http.StatusOK, question)
			})
			
			questions.PUT("/:id", func(c *gin.Context) {
				id := c.Param("id")
				var data map[string]interface{}
				if err := c.ShouldBindJSON(&data); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				
				err := questionRepo.Update(id, data)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully"})
			})
			
			questions.DELETE("/:id", func(c *gin.Context) {
				id := c.Param("id")
				err := questionRepo.Delete(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
			})
			
			// Answers for a question
			questions.GET("/:id/answers", func(c *gin.Context) {
				id := c.Param("id")
				answers, err := answerRepo.GetByQuestionID(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, answers)
			})
			
			questions.POST("/:id/answers", func(c *gin.Context) {
				questionID := c.Param("id")
				var data map[string]interface{}
				if err := c.ShouldBindJSON(&data); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				
				// Make sure the question ID is included
				data["question_id"] = questionID
				
				id, err := answerRepo.Create(data)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusCreated, gin.H{"id": id})
			})
		}

		// Answers endpoints
		answers := api.Group("/answers")
		{
			answers.GET("/:id", func(c *gin.Context) {
				id := c.Param("id")
				answer, err := answerRepo.GetByID(id)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
					return
				}
				
				c.JSON(http.StatusOK, answer)
			})
			
			answers.PUT("/:id", func(c *gin.Context) {
				id := c.Param("id")
				var data map[string]interface{}
				if err := c.ShouldBindJSON(&data); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				
				err := answerRepo.Update(id, data)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{"message": "Answer updated successfully"})
			})
			
			answers.DELETE("/:id", func(c *gin.Context) {
				id := c.Param("id")
				err := answerRepo.Delete(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{"message": "Answer deleted successfully"})
			})
		}
	}

	return router
}

func main() {
	// Initialize database connection
	database := db.GetDB()
	defer database.Close()

	// Setup router with database
	router := setupRouter(database)

	// Start the server
	log.Println("Server is running on port 8080 with Gin")
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
