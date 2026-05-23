package d2script

import (
	"os"
	"testing"
	"time"
)

func TestScriptEngine_EvalTimeout(t *testing.T) {
	s := CreateScriptEngine()
	s.AllowEval()

	// This should timeout
	start := time.Now()
	_, err := s.Eval("while(true) {}")
	elapsed := time.Since(start)

	if err == nil {
		t.Error("expected error for infinite loop, got nil")
	}

	// In release builds, Eval returns ErrScriptingDisabled immediately.
	// We only check the timeout duration for non-release builds where it actually runs.
	if err == ErrEvalTimeout && elapsed < 100*time.Millisecond {
		t.Errorf("timeout too fast: %v", elapsed)
	}
}

func TestScriptEngine_RunScriptTimeout(t *testing.T) {
	s := CreateScriptEngine()
	s.AllowEval()

	// Create a temporary script file with an infinite loop
	tmpFile, err := os.CreateTemp("", "timeout_test_*.js")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	scriptContent := []byte("while(true) {}")
	if _, err := tmpFile.Write(scriptContent); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	start := time.Now()
	_, err = s.RunScript(tmpFile.Name())
	elapsed := time.Since(start)

	if err == nil {
		t.Error("expected error for infinite loop in RunScript, got nil")
	}

	if err == ErrEvalTimeout && elapsed < 100*time.Millisecond {
		t.Errorf("RunScript timeout too fast: %v", elapsed)
	}
}
