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

func TestKappaSystem_ProcessMove(t *testing.T) {
	engine := NewBaalAalEngine()
	event := &IAxiomaticEvent{
		ID:   "move-1",
		Type: "PLAYER_MOVE",
		Metadata: map[string]interface{}{
			"client_id": "player-1",
			"x":         10.5,
			"y":         20.75,
		},
	}

	engine.KappaSystem.HandleMove(event)

	if event.Metadata["kappa_x"] != int32(10500) {
		t.Errorf("expected kappa_x 10500, got %v", event.Metadata["kappa_x"])
	}
	if event.Metadata["kappa_y"] != int32(20750) {
		t.Errorf("expected kappa_y 20750, got %v", event.Metadata["kappa_y"])
	}
}

func TestAxiomaticEventBus_Subscription(t *testing.T) {
	bus := NewAxiomaticEventBus(10)
	called := false
	bus.Subscribe("test", func(e *IAxiomaticEvent) {
		called = true
	})

	bus.Publish(&IAxiomaticEvent{Type: "TEST"})
	if !called {
		t.Error("expected subscriber to be called")
	}

	called = false
	bus.Unsubscribe("test")
	bus.Publish(&IAxiomaticEvent{Type: "TEST"})
	if called {
		t.Error("expected subscriber not to be called after unsubscribe")
	}
}
