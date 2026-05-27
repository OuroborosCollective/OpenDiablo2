package d2emergent

import (
	"testing"
)

func TestMarkgrafARELogik(t *testing.T) {
	// Logic-only test to avoid deep dependencies
	engine := &ARELogikEngine{
		GlobalResonance: 1.0,
		Tick:            0,
		Expansion:       1.0,
		Entropy:         0.5,
		Feedback:        0.0,
	}

	// Run multiple cycles to verify deterministic behavior and cyclical nature
	lastResonance := engine.GlobalResonance
	lastExpansion := engine.Expansion

	for i := 0; i < 100; i++ {
		event := engine.ProcessEmergence()

		if event.Type != "WorldEmergence" {
			t.Errorf("expected event type WorldEmergence, got %s", event.Type)
		}

		resonance := event.Payload.(float64)
		expansion := event.Metadata["expansion"].(float64)
		entropy := event.Metadata["entropy"].(float64)
		tick := event.Metadata["tick"].(uint64)

		if tick != uint64(i+1) {
			t.Errorf("expected tick %d, got %d", i+1, tick)
		}

		// Resonance should change (oscillation)
		if i > 0 && resonance == lastResonance {
			t.Errorf("resonance did not change at tick %d", tick)
		}

		// Expansion should be monotonically increasing (until reset)
		if expansion <= lastExpansion && expansion != 1.0 {
			t.Errorf("expansion did not increase at tick %d: %f -> %f", tick, lastExpansion, expansion)
		}

		// Entropy should be within [0, 2.0]
		if entropy < 0 || entropy > 2.0 {
			t.Errorf("entropy out of bounds at tick %d: %f", tick, entropy)
		}

		lastResonance = resonance
		lastExpansion = expansion
	}
}

func TestDeterministicMarkgraf(t *testing.T) {
	engine1 := &ARELogikEngine{
		GlobalResonance: 1.0,
		Tick:            0,
		Expansion:       1.0,
		Entropy:         0.5,
		Feedback:        0.0,
	}
	engine2 := &ARELogikEngine{
		GlobalResonance: 1.0,
		Tick:            0,
		Expansion:       1.0,
		Entropy:         0.5,
		Feedback:        0.0,
	}

	for i := 0; i < 50; i++ {
		event1 := engine1.ProcessEmergence()
		event2 := engine2.ProcessEmergence()

		res1 := event1.Payload.(float64)
		res2 := event2.Payload.(float64)

		if res1 != res2 {
			t.Errorf("non-deterministic resonance at tick %d: %f != %f", i+1, res1, res2)
		}
	}
}
