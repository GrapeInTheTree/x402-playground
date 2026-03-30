package quiz

import "sync"

// Question defines a coding quiz question.
type Question struct {
	ID          string
	Title       string
	Description string // shown in TUI before opening editor
	Difficulty  string // "easy", "medium", "hard"
	Category    string // "Basics", "ERC-20", "EIP-712", "EIP-3009", "EIP-2612", "x402", "Permit2"
	Language    Lang   // LangGo or LangSolidity
	Template    string // source code with TODO sections
	TestCode    string // test code to validate the solution
	Hints       []string
}

// Result holds the outcome of running tests on a submitted solution.
type Result struct {
	Compiled bool
	Passed   int
	Total    int
	Output   string // raw test output
	Error    string // compilation or runtime error
}

// Score tracks quiz progress.
type Score struct {
	Answered  int
	Correct   int
	Questions int
}

// ModuleProgress tracks progress for a single question group.
type ModuleProgress struct {
	Name      string
	Total     int
	Attempted int
	Passed    int
}

// QuizProgress tracks overall quiz progress shared between Learn and Dashboard.
// Thread-safe: use getter/setter methods for concurrent access.
type QuizProgress struct {
	mu      sync.RWMutex
	Modules []ModuleProgress
	Score   Score
}

// SetModules updates the module progress list.
func (qp *QuizProgress) SetModules(modules []ModuleProgress) {
	qp.mu.Lock()
	defer qp.mu.Unlock()
	qp.Modules = modules
}

// GetModules returns a copy of the module progress list.
func (qp *QuizProgress) GetModules() []ModuleProgress {
	qp.mu.RLock()
	defer qp.mu.RUnlock()
	out := make([]ModuleProgress, len(qp.Modules))
	copy(out, qp.Modules)
	return out
}

// SetScore updates the score.
func (qp *QuizProgress) SetScore(s Score) {
	qp.mu.Lock()
	defer qp.mu.Unlock()
	qp.Score = s
}

// GetScore returns the current score.
func (qp *QuizProgress) GetScore() Score {
	qp.mu.RLock()
	defer qp.mu.RUnlock()
	return qp.Score
}
