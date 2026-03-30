package quiz

import (
	"sync"
	"testing"
)

func TestQuizProgress_ThreadSafety(t *testing.T) {
	qp := &QuizProgress{}

	var wg sync.WaitGroup

	// Writer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			qp.SetModules([]ModuleProgress{
				{Name: "test", Total: 10, Attempted: i, Passed: i},
			})
			qp.SetScore(Score{Answered: i, Correct: i, Questions: 10})
		}
	}()

	// Reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			modules := qp.GetModules()
			_ = modules
			score := qp.GetScore()
			_ = score
		}
	}()

	wg.Wait()
}

func TestQuizProgress_GetModulesReturnsCopy(t *testing.T) {
	qp := &QuizProgress{}
	qp.SetModules([]ModuleProgress{
		{Name: "Go Basics", Total: 5, Passed: 3},
	})

	modules := qp.GetModules()
	if len(modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(modules))
	}

	// Mutating the returned copy should not affect the original
	modules[0].Passed = 999
	original := qp.GetModules()
	if original[0].Passed != 3 {
		t.Errorf("GetModules should return a copy, but original was mutated to %d", original[0].Passed)
	}
}

func TestQuizProgress_SetAndGetScore(t *testing.T) {
	qp := &QuizProgress{}
	qp.SetScore(Score{Answered: 5, Correct: 3, Questions: 10})

	score := qp.GetScore()
	if score.Answered != 5 || score.Correct != 3 || score.Questions != 10 {
		t.Errorf("unexpected score: %+v", score)
	}
}
