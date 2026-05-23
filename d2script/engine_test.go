package d2script

import (
	"testing"
	"time"
)

func TestScriptEngine_EvalTimeout(t *testing.T) {
	// Skip this test in release builds as scripting is disabled anyway
	if isRelease() {
		t.Skip("skipping test in release build")
	}

	s := CreateScriptEngine()
	s.AllowEval()

	// This should timeout
	start := time.Now()
	_, err := s.Eval("while(true) {}")
	elapsed := time.Since(start)

	if err == nil {
		t.Error("expected error for infinite loop, got nil")
	}

	if elapsed < 100*time.Millisecond {
		t.Errorf("timeout too fast: %v", elapsed)
	}
}

func isRelease() bool {
	s := CreateScriptEngine()
	_, err := s.Eval("1+1")
	return err != nil && err.Error() == "arbitrary scripting disabled in release builds"
}
