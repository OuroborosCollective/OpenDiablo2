package d2script

import (
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
