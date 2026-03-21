package quiz

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
type QuizProgress struct {
	Modules []ModuleProgress
	Score   Score
}
