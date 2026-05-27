package d2script

import (
	"testing"
)

func TestAxiomaticSystem_Integration(t *testing.T) {
	engine := NewBaalAalEngine()
	worldSys := NewWorldSystem()

	// Register world system rule
	engine.RegisterRule("WorldEmergence", worldSys.HandleEmergence)

	// 1. Simulate Player Move
	moveEvent := &IAxiomaticEvent{
		ID:   "move-integration",
		Type: "PlayerMove",
		Payload: map[string]interface{}{
			"x": 100.0,
			"y": 200.0,
		},
		Metadata: map[string]interface{}{
			"client_id": "test-client",
			"x":         100.0,
			"y":         200.0,
		},
	}

	engine.EventBus.Publish(moveEvent)

	// Process cycle to trigger rules (though KappaSystem currently processes on Publish via subscription)
	engine.ProcessCycle(1)

	// Verify KappaSystem update
	engine.KappaSystem.RLock()
	pos, ok := engine.KappaSystem.Positions["test-client"]
	engine.KappaSystem.RUnlock()

	if !ok {
		t.Fatal("expected position for test-client in KappaSystem")
	}
	if pos[0] != 100000 || pos[1] != 200000 {
		t.Errorf("expected [100000, 200000], got %v", pos)
	}

	// 2. Simulate World Emergence
	emergenceEvent := &IAxiomaticEvent{
		ID:      "emergence-integration",
		Type:    "WorldEmergence",
		Payload: 123.456,
	}

	engine.EventBus.Publish(emergenceEvent)

	// Process cycle to trigger WorldSystem.HandleEmergence
	engine.ProcessCycle(2)

	if worldSys.GlobalResonance != 123.456 {
		t.Errorf("expected GlobalResonance 123.456, got %f", worldSys.GlobalResonance)
	}

	// 3. Verify Resonance and Cycle Status
	resonance, cycle := engine.GetStatus()
	if resonance < 0 || resonance >= 1.0 {
		t.Errorf("invalid resonance: %f", resonance)
	}
	if cycle <= 0 {
		t.Errorf("invalid cycle: %f", cycle)
	}
}
