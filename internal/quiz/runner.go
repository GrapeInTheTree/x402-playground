package quiz

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Lang identifies the quiz language.
type Lang string

const (
	LangGo       Lang = "go"
	LangSolidity Lang = "solidity"
)

// Runner executes quiz solutions in a temporary project.
type Runner struct {
	workDir string
	lang    Lang
}

// NewRunner creates a temp directory for Go quiz solutions.
func NewRunner() (*Runner, error) {
	return newRunnerForLang(LangGo)
}

// NewSolidityRunner creates a temp Foundry project for Solidity quiz solutions.
func NewSolidityRunner() (*Runner, error) {
	return newRunnerForLang(LangSolidity)
}

func newRunnerForLang(lang Lang) (*Runner, error) {
	dir, err := os.MkdirTemp("", "x402-quiz-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}

	switch lang {
	case LangGo:
		modContent := "module x402quiz\n\ngo 1.21\n"
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644); err != nil {
			os.RemoveAll(dir)
			return nil, fmt.Errorf("write go.mod: %w", err)
		}
	case LangSolidity:
		if err := initFoundryProject(dir); err != nil {
			os.RemoveAll(dir)
			return nil, err
		}
	}

	return &Runner{workDir: dir, lang: lang}, nil
}

func initFoundryProject(dir string) error {
	// Create directory structure
	for _, sub := range []string{"src", "test"} {
		if err := os.MkdirAll(filepath.Join(dir, sub), 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", sub, err)
		}
	}

	// foundry.toml
	toml := `[profile.default]
src = "src"
out = "out"
libs = ["lib"]
solc_version = "0.8.20"
`
	if err := os.WriteFile(filepath.Join(dir, "foundry.toml"), []byte(toml), 0644); err != nil {
		return fmt.Errorf("write foundry.toml: %w", err)
	}

	// Initialize git repo (forge install requires it)
	gitInit := exec.Command("git", "init")
	gitInit.Dir = dir
	if out, err := gitInit.CombinedOutput(); err != nil {
		return fmt.Errorf("git init: %w\n%s", err, out)
	}

	// Install forge-std for testing
	cmd := exec.Command("forge", "install", "foundry-rs/forge-std", "--no-commit")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("forge install forge-std: %w\n%s", err, out)
	}

	// remappings.txt
	remappings := "forge-std/=lib/forge-std/src/\n"
	if err := os.WriteFile(filepath.Join(dir, "remappings.txt"), []byte(remappings), 0644); err != nil {
		return fmt.Errorf("write remappings.txt: %w", err)
	}

	return nil
}

// Run writes the solution and test files, then runs the appropriate test command.
func (r *Runner) Run(solution, testCode string) *Result {
	switch r.lang {
	case LangGo:
		return r.runGo(solution, testCode)
	case LangSolidity:
		return r.runSolidity(solution, testCode)
	default:
		return &Result{Error: fmt.Sprintf("unknown language: %s", r.lang)}
	}
}

func (r *Runner) runGo(solution, testCode string) *Result {
	solPath := filepath.Join(r.workDir, "solution.go")
	testPath := filepath.Join(r.workDir, "solution_test.go")

	if err := os.WriteFile(solPath, []byte(solution), 0644); err != nil {
		return &Result{Error: fmt.Sprintf("write solution: %v", err)}
	}
	if err := os.WriteFile(testPath, []byte(testCode), 0644); err != nil {
		return &Result{Error: fmt.Sprintf("write test: %v", err)}
	}

	cmd := exec.Command("go", "test", "-v", "-count=1", "./...")
	cmd.Dir = r.workDir
	out, err := cmd.CombinedOutput()
	output := string(out)

	result := &Result{Output: output}

	if err != nil {
		if strings.Contains(output, "build failed") ||
			strings.Contains(output, "cannot") ||
			strings.Contains(output, "undefined") ||
			strings.Contains(output, "syntax error") {
			result.Error = "Compilation failed"
			return result
		}
		result.Compiled = true
	} else {
		result.Compiled = true
	}

	result.Passed = strings.Count(output, "--- PASS")
	result.Total = result.Passed + strings.Count(output, "--- FAIL")

	if result.Total == 0 && result.Compiled {
		result.Total = 1
		if err == nil {
			result.Passed = 1
		}
	}

	return result
}

func (r *Runner) runSolidity(solution, testCode string) *Result {
	solPath := filepath.Join(r.workDir, "src", "Solution.sol")
	testPath := filepath.Join(r.workDir, "test", "Solution.t.sol")

	if err := os.WriteFile(solPath, []byte(solution), 0644); err != nil {
		return &Result{Error: fmt.Sprintf("write solution: %v", err)}
	}
	if err := os.WriteFile(testPath, []byte(testCode), 0644); err != nil {
		return &Result{Error: fmt.Sprintf("write test: %v", err)}
	}

	cmd := exec.Command("forge", "test", "-vv")
	cmd.Dir = r.workDir
	out, err := cmd.CombinedOutput()
	output := string(out)

	result := &Result{Output: output}

	if err != nil {
		if strings.Contains(output, "Compiler run failed") ||
			strings.Contains(output, "Error") && strings.Contains(output, "-->") {
			result.Error = "Compilation failed"
			return result
		}
		result.Compiled = true
	} else {
		result.Compiled = true
	}

	// Parse forge test output: [PASS] or [FAIL]
	result.Passed = strings.Count(output, "[PASS]")
	result.Total = result.Passed + strings.Count(output, "[FAIL]")

	if result.Total == 0 && result.Compiled {
		result.Total = 1
		if err == nil {
			result.Passed = 1
		}
	}

	return result
}

// TemplatePath returns the path where the solution file is written.
func (r *Runner) TemplatePath() string {
	switch r.lang {
	case LangSolidity:
		return filepath.Join(r.workDir, "src", "Solution.sol")
	default:
		return filepath.Join(r.workDir, "solution.go")
	}
}

// Cleanup removes the temporary directory.
func (r *Runner) Cleanup() {
	os.RemoveAll(r.workDir)
}
