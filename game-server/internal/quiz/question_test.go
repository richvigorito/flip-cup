package quiz

import "testing"

func TestCheckAnswer_ExactMatch(t *testing.T) {
	q := &Question{
		Prompt:  "What is the capital of France?",
		Answers: []string{"Paris"},
	}
	if !q.CheckAnswer("Paris") {
		t.Error("Expected exact match to return true")
	}
}

func TestCheckAnswer_CaseInsensitive(t *testing.T) {
	q := &Question{
		Prompt:  "What is the capital of France?",
		Answers: []string{"Paris"},
	}
	if !q.CheckAnswer("paris") {
		t.Error("Expected case-insensitive match to return true")
	}
	if !q.CheckAnswer("PARIS") {
		t.Error("Expected uppercase match to return true")
	}
}

func TestCheckAnswer_TrimsWhitespace(t *testing.T) {
	q := &Question{
		Prompt:  "What is the capital of France?",
		Answers: []string{"Paris"},
	}
	if !q.CheckAnswer("  Paris  ") {
		t.Error("Expected answer with leading/trailing spaces to return true")
	}
}

func TestCheckAnswer_RemovesPunctuation(t *testing.T) {
	q := &Question{
		Prompt:  "Who wrote Hamlet?",
		Answers: []string{"Shakespeare"},
	}
	if !q.CheckAnswer("Shakespeare!") {
		t.Error("Expected answer with punctuation to return true")
	}
	if !q.CheckAnswer("Shake-speare") {
		t.Error("Expected answer with hyphen to return true")
	}
}

func TestCheckAnswer_IncorrectAnswer(t *testing.T) {
	q := &Question{
		Prompt:  "What is the capital of France?",
		Answers: []string{"Paris"},
	}
	if q.CheckAnswer("London") {
		t.Error("Expected wrong answer to return false")
	}
}

func TestCheckAnswer_EmptyInput(t *testing.T) {
	q := &Question{
		Prompt:  "What is the capital of France?",
		Answers: []string{"Paris"},
	}
	if q.CheckAnswer("") {
		t.Error("Expected empty input to return false")
	}
}

func TestCheckAnswer_MultipleCorrectAnswers(t *testing.T) {
	q := &Question{
		Prompt:  "Name a primary color.",
		Answers: []string{"Red", "Blue", "Yellow"},
	}
	if !q.CheckAnswer("red") {
		t.Error("Expected first answer to match")
	}
	if !q.CheckAnswer("Blue") {
		t.Error("Expected second answer to match")
	}
	if !q.CheckAnswer("YELLOW") {
		t.Error("Expected third answer to match")
	}
}

func TestCheckAnswer_SpecialCharactersInAnswer(t *testing.T) {
	q := &Question{
		Prompt:  "What is the answer?",
		Answers: []string{"St. Louis"},
	}
	// normalize strips the dot and space, so "St. Louis" -> "stlouis"
	if !q.CheckAnswer("St Louis") {
		t.Error("Expected answer without punctuation to match normalized correct answer")
	}
	if !q.CheckAnswer("st. louis") {
		t.Error("Expected lowercase with period to match")
	}
}

func TestCheckAnswer_NumbersInAnswer(t *testing.T) {
	q := &Question{
		Prompt:  "How many sides does a hexagon have?",
		Answers: []string{"6"},
	}
	if !q.CheckAnswer("6") {
		t.Error("Expected numeric answer to match")
	}
	if q.CheckAnswer("7") {
		t.Error("Expected wrong number to return false")
	}
}
