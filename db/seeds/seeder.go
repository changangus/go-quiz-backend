package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/changangus/go-quiz-backend/internal/models"
	"github.com/changangus/go-quiz-backend/internal/repository"

	_ "github.com/lib/pq"
)

type (
	Quiz     = models.Quiz
	Question = models.Question
	Answer   = models.Answer
)

var (
	// CreateQuiz creates a new quiz in the database
	createQuiz = repository.CreateQuiz
	// CreateQuestion creates a new question in the database
	createQuestion = repository.CreateQuestion
	// CreateAnswer creates a new answer in the database
	createAnswer = repository.CreateAnswer
)

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Create tables if they don't exist
func alertUserIfTablesDontExist(db *sql.DB) error {
	// Check if tables exist
	var tableExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'quizzes')").Scan(&tableExists)
	if err != nil {
		return err
	}

	if !tableExists {
		fmt.Println("Tables don't exist, run migrations")
	}

	return nil
}

func main() {
	// Get environment variables or use defaults
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "quizdb")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")

	// PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database")

	// First, check if the tables exist and create them if they don't
	err = alertUserIfTablesDontExist(db)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Create a WCAG 2.2 Section 1 quiz
	quizID, err := createQuiz(db, &Quiz{
		Title:       "WCAG 2.2 Section 1: Perceivable",
		Description: "Test your knowledge of WCAG 2.2 Section 1 guidelines on making information and user interface components perceivable to users.",
	})
	if err != nil {
		log.Fatal("Failed to create quiz:", err)
	}

	// Seed questions and answers
	seedQuestions(db, quizID)

	fmt.Println("Database seeded successfully!")
}

func seedQuestions(db *sql.DB, quizID int64) {
	// Guideline 1.1: Text Alternatives
	q1ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 guideline 1.1.1 (Non-text Content), what must be provided for all non-text content?",
		Type:     "multiple_choice",
		Order:    1,
	})
	if err != nil {
		log.Fatal("Failed to create question 1:", err)
	}

	answers := []Answer{
		{QuestionID: int(q1ID), Answer: "A text alternative that serves the equivalent purpose", Correct: true},
		{QuestionID: int(q1ID), Answer: "A detailed image description regardless of purpose", Correct: false},
		{QuestionID: int(q1ID), Answer: "An audio file explaining the content", Correct: false},
		{QuestionID: int(q1ID), Answer: "A simplified version of the content", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Guideline 1.2: Time-based Media
	q2ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "For prerecorded audio-only content, what is required to meet WCAG 2.2 Success Criterion 1.2.1 (Audio-only and Video-only)?",
		Type:     "multiple_choice",
		Order:    2,
	})
	if err != nil {
		log.Fatal("Failed to create question 2:", err)
	}

	answers = []Answer{
		{QuestionID: int(q2ID), Answer: "An alternative for time-based media that presents equivalent information", Correct: true},
		{QuestionID: int(q2ID), Answer: "A sign language interpretation", Correct: false},
		{QuestionID: int(q2ID), Answer: "Background music", Correct: false},
		{QuestionID: int(q2ID), Answer: "A volume control", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// More on Time-based Media
	q3ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "What is required for prerecorded video content to meet WCAG 2.2 Success Criterion 1.2.2 (Captions)?",
		Type:     "multiple_choice",
		Order:    3,
	})
	if err != nil {
		log.Fatal("Failed to create question 3:", err)
	}

	answers = []Answer{
		{QuestionID: int(q3ID), Answer: "Captions for all prerecorded audio content", Correct: true},
		{QuestionID: int(q3ID), Answer: "An audio description track", Correct: false},
		{QuestionID: int(q3ID), Answer: "A sign language interpretation", Correct: false},
		{QuestionID: int(q3ID), Answer: "A full text transcript only", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Audio Description
	q4ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.2.3, what must be provided for prerecorded video content?",
		Type:     "multiple_choice",
		Order:    4,
	})
	if err != nil {
		log.Fatal("Failed to create question 4:", err)
	}

	answers = []Answer{
		{QuestionID: int(q4ID), Answer: "An audio description or media alternative", Correct: true},
		{QuestionID: int(q4ID), Answer: "Only captions", Correct: false},
		{QuestionID: int(q4ID), Answer: "Only a transcript", Correct: false},
		{QuestionID: int(q4ID), Answer: "Only sign language", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Guideline 1.3: Adaptable
	q5ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 guideline 1.3.1 (Info and Relationships), what should be programmatically determined or available in text?",
		Type:     "multiple_choice",
		Order:    5,
	})
	if err != nil {
		log.Fatal("Failed to create question 5:", err)
	}

	answers = []Answer{
		{QuestionID: int(q5ID), Answer: "Information, structure, and relationships conveyed through presentation", Correct: true},
		{QuestionID: int(q5ID), Answer: "Only color-based information", Correct: false},
		{QuestionID: int(q5ID), Answer: "Only heading structures", Correct: false},
		{QuestionID: int(q5ID), Answer: "Only form labels", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Meaningful Sequence
	q6ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "What does WCAG 2.2 Success Criterion 1.3.2 (Meaningful Sequence) require?",
		Type:     "multiple_choice",
		Order:    6,
	})
	if err != nil {
		log.Fatal("Failed to create question 6:", err)
	}

	answers = []Answer{
		{QuestionID: int(q6ID), Answer: "When the sequence of content affects its meaning, a correct reading sequence can be programmatically determined", Correct: true},
		{QuestionID: int(q6ID), Answer: "Content must always be presented in the same sequence", Correct: false},
		{QuestionID: int(q6ID), Answer: "Users must be able to rearrange content in any sequence", Correct: false},
		{QuestionID: int(q6ID), Answer: "All content must be organized in alphabetical order", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Sensory Characteristics
	q7ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.3.3 (Sensory Characteristics), instructions for understanding content should not rely solely on what?",
		Type:     "multiple_choice",
		Order:    7,
	})
	if err != nil {
		log.Fatal("Failed to create question 7:", err)
	}

	answers = []Answer{
		{QuestionID: int(q7ID), Answer: "Sensory characteristics such as shape, color, size, visual location, orientation, or sound", Correct: true},
		{QuestionID: int(q7ID), Answer: "Text-based instructions", Correct: false},
		{QuestionID: int(q7ID), Answer: "Keyboard shortcuts", Correct: false},
		{QuestionID: int(q7ID), Answer: "Menu selections", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Orientation (1.3.4)
	q8ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "What does WCAG 2.2 Success Criterion 1.3.4 (Orientation) require?",
		Type:     "multiple_choice",
		Order:    8,
	})
	if err != nil {
		log.Fatal("Failed to create question 8:", err)
	}

	answers = []Answer{
		{QuestionID: int(q8ID), Answer: "Content does not restrict its view and operation to a single display orientation, unless a specific orientation is essential", Correct: true},
		{QuestionID: int(q8ID), Answer: "All content must work in landscape mode only", Correct: false},
		{QuestionID: int(q8ID), Answer: "All content must work in portrait mode only", Correct: false},
		{QuestionID: int(q8ID), Answer: "Users must manually select their preferred orientation", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Identify Input Purpose (1.3.5)
	q9ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.3.5 (Identify Input Purpose), what should be true about input fields that collect information about the user?",
		Type:     "multiple_choice",
		Order:    9,
	})
	if err != nil {
		log.Fatal("Failed to create question 9:", err)
	}

	answers = []Answer{
		{QuestionID: int(q9ID), Answer: "The purpose of each input field can be programmatically determined", Correct: true},
		{QuestionID: int(q9ID), Answer: "All input fields must be optional", Correct: false},
		{QuestionID: int(q9ID), Answer: "Input fields should not collect personal information", Correct: false},
		{QuestionID: int(q9ID), Answer: "Input fields must always use autocomplete", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Guideline 1.4: Distinguishable
	q10ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.4.1 (Use of Color), color should not be used as what?",
		Type:     "multiple_choice",
		Order:    10,
	})
	if err != nil {
		log.Fatal("Failed to create question 10:", err)
	}

	answers = []Answer{
		{QuestionID: int(q10ID), Answer: "The only visual means of conveying information, indicating an action, prompting a response, or distinguishing a visual element", Correct: true},
		{QuestionID: int(q10ID), Answer: "A decorative element", Correct: false},
		{QuestionID: int(q10ID), Answer: "A way to highlight text", Correct: false},
		{QuestionID: int(q10ID), Answer: "A branding element", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Audio Control
	q11ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "What does WCAG 2.2 Success Criterion 1.4.2 (Audio Control) require for any audio that plays automatically for more than 3 seconds?",
		Type:     "multiple_choice",
		Order:    11,
	})
	if err != nil {
		log.Fatal("Failed to create question 11:", err)
	}

	answers = []Answer{
		{QuestionID: int(q11ID), Answer: "A mechanism to pause or stop the audio, or a mechanism to control audio volume independently from the overall system volume", Correct: true},
		{QuestionID: int(q11ID), Answer: "Audio must never play automatically", Correct: false},
		{QuestionID: int(q11ID), Answer: "Audio must stop automatically after 10 seconds", Correct: false},
		{QuestionID: int(q11ID), Answer: "Audio must always include captions", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Contrast Minimum
	q12ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "What is the minimum contrast ratio required for normal text according to WCAG 2.2 Success Criterion 1.4.3 (Contrast Minimum)?",
		Type:     "multiple_choice",
		Order:    12,
	})
	if err != nil {
		log.Fatal("Failed to create question 12:", err)
	}

	answers = []Answer{
		{QuestionID: int(q12ID), Answer: "4.5:1", Correct: true},
		{QuestionID: int(q12ID), Answer: "3:1", Correct: false},
		{QuestionID: int(q12ID), Answer: "7:1", Correct: false},
		{QuestionID: int(q12ID), Answer: "2:1", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Resize Text
	q13ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.4.4 (Resize Text), text should be able to be resized without assistive technology up to what percentage without loss of content or functionality?",
		Type:     "multiple_choice",
		Order:    13,
	})
	if err != nil {
		log.Fatal("Failed to create question 13:", err)
	}

	answers = []Answer{
		{QuestionID: int(q13ID), Answer: "200%", Correct: true},
		{QuestionID: int(q13ID), Answer: "150%", Correct: false},
		{QuestionID: int(q13ID), Answer: "300%", Correct: false},
		{QuestionID: int(q13ID), Answer: "100%", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Images of Text
	q14ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.4.5 (Images of Text), when should text be used instead of images of text?",
		Type:     "multiple_choice",
		Order:    14,
	})
	if err != nil {
		log.Fatal("Failed to create question 14:", err)
	}

	answers = []Answer{
		{QuestionID: int(q14ID), Answer: "Whenever possible, except when a particular presentation of text is essential to the information being conveyed", Correct: true},
		{QuestionID: int(q14ID), Answer: "Only when the text is longer than 20 words", Correct: false},
		{QuestionID: int(q14ID), Answer: "Only when the images cannot be resized", Correct: false},
		{QuestionID: int(q14ID), Answer: "Only when high contrast is required", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Reflow
	q15ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.4.10 (Reflow), content should be presentable without loss of information or functionality at what viewport width?",
		Type:     "multiple_choice",
		Order:    15,
	})
	if err != nil {
		log.Fatal("Failed to create question 15:", err)
	}

	answers = []Answer{
		{QuestionID: int(q15ID), Answer: "320 CSS pixels", Correct: true},
		{QuestionID: int(q15ID), Answer: "240 CSS pixels", Correct: false},
		{QuestionID: int(q15ID), Answer: "480 CSS pixels", Correct: false},
		{QuestionID: int(q15ID), Answer: "640 CSS pixels", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Non-Text Contrast
	q16ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "What minimum contrast ratio is required for user interface components and graphical objects according to WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast)?",
		Type:     "multiple_choice",
		Order:    16,
	})
	if err != nil {
		log.Fatal("Failed to create question 16:", err)
	}

	answers = []Answer{
		{QuestionID: int(q16ID), Answer: "3:1", Correct: true},
		{QuestionID: int(q16ID), Answer: "4.5:1", Correct: false},
		{QuestionID: int(q16ID), Answer: "2:1", Correct: false},
		{QuestionID: int(q16ID), Answer: "7:1", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Text Spacing
	q17ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.4.12 (Text Spacing), no loss of content or functionality should occur when users modify which of the following text properties?",
		Type:     "multiple_choice",
		Order:    17,
	})
	if err != nil {
		log.Fatal("Failed to create question 17:", err)
	}

	answers = []Answer{
		{QuestionID: int(q17ID), Answer: "Line height, spacing between paragraphs, letter spacing, and word spacing", Correct: true},
		{QuestionID: int(q17ID), Answer: "Only font size and color", Correct: false},
		{QuestionID: int(q17ID), Answer: "Only text alignment and indentation", Correct: false},
		{QuestionID: int(q17ID), Answer: "Only font family and style", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// Content on Hover or Focus
	q18ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2 Success Criterion 1.4.13 (Content on Hover or Focus), which requirement applies to additional content that appears on hover or focus?",
		Type:     "multiple_choice",
		Order:    18,
	})
	if err != nil {
		log.Fatal("Failed to create question 18:", err)
	}

	answers = []Answer{
		{QuestionID: int(q18ID), Answer: "It must be dismissable, hoverable, and persistent", Correct: true},
		{QuestionID: int(q18ID), Answer: "It must always disappear after 3 seconds", Correct: false},
		{QuestionID: int(q18ID), Answer: "It must always appear in the top-right corner", Correct: false},
		{QuestionID: int(q18ID), Answer: "It must always include an icon", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	// True/False questions
	q19ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2, content should rely solely on color to convey important information.",
		Type:     "true_false",
		Order:    19,
	})
	if err != nil {
		log.Fatal("Failed to create question 19:", err)
	}

	answers = []Answer{
		{QuestionID: int(q19ID), Answer: "True", Correct: false},
		{QuestionID: int(q19ID), Answer: "False", Correct: true},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}

	q20ID, err := createQuestion(db, &Question{
		QuizID:   int(quizID),
		Question: "According to WCAG 2.2, all audio that plays automatically for more than 3 seconds must have a mechanism to pause, stop, or control the volume independently from the system volume.",
		Type:     "true_false",
		Order:    20,
	})
	if err != nil {
		log.Fatal("Failed to create question 20:", err)
	}

	answers = []Answer{
		{QuestionID: int(q20ID), Answer: "True", Correct: true},
		{QuestionID: int(q20ID), Answer: "False", Correct: false},
	}
	for _, a := range answers {
		if err := createAnswer(db, &a); err != nil {
			log.Fatal("Failed to create answer:", err)
		}
	}
}
