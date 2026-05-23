package d2emergent

import (
	"testing"
)

func TestARELogikEngine_ProcessEmergence(t *testing.T) {
	e := CreateARELogikEngine(5) // LogLevelInfo
	initialResonance := e.GlobalResonance

	e.ProcessEmergence()

	if e.GlobalResonance <= initialResonance {
		t.Errorf("expected resonance to increase, got initial %.2f and new %.2f", initialResonance, e.GlobalResonance)
	}
}
