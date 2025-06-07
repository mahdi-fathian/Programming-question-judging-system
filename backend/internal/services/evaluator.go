package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/onlinejudge/backend/internal/models"
	"github.com/onlinejudge/backend/pkg/database"
)

type EvaluationResult struct {
	Status     string
	TimeUsed   int
	MemoryUsed int
	Error      string
}

type Evaluator struct {
	workDir string
}

func NewEvaluator() *Evaluator {
	workDir := filepath.Join(os.TempDir(), "onlinejudge")
	os.MkdirAll(workDir, 0755)
	return &Evaluator{workDir: workDir}
}

func (e *Evaluator) Evaluate(submission *models.Submission) error {
	// Create a unique directory for this submission
	submissionDir := filepath.Join(e.workDir, fmt.Sprintf("submission_%d", submission.ID))
	os.MkdirAll(submissionDir, 0755)
	defer os.RemoveAll(submissionDir)

	// Write code to file
	codeFile := filepath.Join(submissionDir, getSourceFileName(submission.Language))
	if err := os.WriteFile(codeFile, []byte(submission.Code), 0644); err != nil {
		return fmt.Errorf("failed to write code file: %v", err)
	}

	// Compile if needed
	if err := e.compile(submission.Language, codeFile); err != nil {
		submission.Status = "compilation_error"
		submission.Error = err.Error()
		database.DB.Save(submission)
		return nil
	}

	// Get test cases
	var testCases []models.TestCase
	database.DB.Where("problem_id = ?", submission.ProblemID).Find(&testCases)

	// Run test cases
	var results []models.SubmissionResult
	for _, tc := range testCases {
		result := e.runTestCase(submission, tc, codeFile)
		results = append(results, models.SubmissionResult{
			SubmissionID: submission.ID,
			TestCaseID:   tc.ID,
			Status:       result.Status,
			TimeUsed:     result.TimeUsed,
			MemoryUsed:   result.MemoryUsed,
			Error:        result.Error,
		})

		// If any test case fails, stop evaluation
		if result.Status != "accepted" {
			break
		}
	}

	// Save results
	for _, result := range results {
		database.DB.Create(&result)
	}

	// Update submission status
	submission.Status = results[len(results)-1].Status
	submission.TimeUsed = results[len(results)-1].TimeUsed
	submission.MemoryUsed = results[len(results)-1].MemoryUsed
	database.DB.Save(submission)

	return nil
}

func (e *Evaluator) compile(language, codeFile string) error {
	switch language {
	case "cpp":
		cmd := exec.Command("g++", "-std=c++17", "-O2", codeFile, "-o", codeFile+".out")
		return cmd.Run()
	case "java":
		cmd := exec.Command("javac", codeFile)
		return cmd.Run()
	case "python":
		// Python is interpreted, no compilation needed
		return nil
	default:
		return fmt.Errorf("unsupported language: %s", language)
	}
}

func (e *Evaluator) runTestCase(submission *models.Submission, tc models.TestCase, codeFile string) EvaluationResult {
	var cmd *exec.Cmd
	switch submission.Language {
	case "cpp":
		cmd = exec.Command(codeFile + ".out")
	case "java":
		cmd = exec.Command("java", "-cp", filepath.Dir(codeFile), filepath.Base(codeFile[:len(codeFile)-5]))
	case "python":
		cmd = exec.Command("python", codeFile)
	default:
		return EvaluationResult{Status: "runtime_error", Error: "Unsupported language"}
	}

	// Set up input/output
	cmd.Stdin = bytes.NewReader([]byte(tc.Input))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set resource limits
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Start the process
	startTime := time.Now()
	err := cmd.Start()
	if err != nil {
		return EvaluationResult{Status: "runtime_error", Error: err.Error()}
	}

	// Create a channel to receive the result
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	// Set timeout
	timeout := time.Duration(submission.Problem.TimeLimit) * time.Millisecond
	select {
	case err := <-done:
		timeUsed := int(time.Since(startTime).Milliseconds())
		if err != nil {
			return EvaluationResult{
				Status:   "runtime_error",
				TimeUsed: timeUsed,
				Error:    stderr.String(),
			}
		}

		// Check output
		output := stdout.String()
		if output != tc.Output {
			return EvaluationResult{
				Status:   "wrong_answer",
				TimeUsed: timeUsed,
				Error:    "Output does not match expected output",
			}
		}

		return EvaluationResult{
			Status:   "accepted",
			TimeUsed: timeUsed,
		}

	case <-time.After(timeout):
		// Kill the process group
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return EvaluationResult{
			Status:   "time_limit",
			TimeUsed: submission.Problem.TimeLimit,
			Error:    "Time limit exceeded",
		}
	}
}

func getSourceFileName(language string) string {
	switch language {
	case "cpp":
		return "main.cpp"
	case "java":
		return "Main.java"
	case "python":
		return "main.py"
	default:
		return "main"
	}
} 