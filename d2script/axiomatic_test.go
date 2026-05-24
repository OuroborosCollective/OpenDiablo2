package d2script

import (
	"testing"
)

func TestAREStateCompiler_ToKappa(t *testing.T) {
	c := &AREStateCompiler{}
	if c.ToKappa(1.2345) != 1234 {
		t.Errorf("expected 1234, got %d", c.ToKappa(1.2345))
	}
}

func TestAxiomaticEventBus_Publish(t *testing.T) {
	bus := NewAxiomaticEventBus(10)
	event := &IAxiomaticEvent{
		ID:   "test-1",
		Type: "PLAYER_MOVE",
	}

	bus.Publish(event)

	if event.SequenceID != 0 {
		t.Errorf("expected sequence ID 0, got %d", event.SequenceID)
	}

	if event.Metadata["baalaal_cycle"] == nil {
		t.Error("expected baalaal_cycle in metadata")
	}

	history := bus.GetHistory()
	if len(history) != 1 {
		t.Errorf("expected history length 1, got %d", len(history))
	}
}

func TestBaalAalEngine_ProcessCycle(t *testing.T) {
	engine := NewBaalAalEngine()
	initialResonance := engine.EventBus.resonanceState
	engine.ProcessCycle(1)
	if engine.EventBus.resonanceState == initialResonance {
		t.Error("expected resonance state to change after cycle")
	}
}

func TestWorldSystem_HandleEmergence(t *testing.T) {
	ws := NewWorldSystem()
	event := &IAxiomaticEvent{
		Payload: 123.456,
	}

	ws.HandleEmergence(event)

	if ws.GlobalResonance != 123.456 {
		t.Errorf("expected resonance 123.456, got %.3f", ws.GlobalResonance)
	}
}
