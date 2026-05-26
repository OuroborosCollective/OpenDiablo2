package d2script

import (
	"testing"
)

func TestAxiomaticEventBus_Subscribe(t *testing.T) {
	bus := NewAxiomaticEventBus(10)
	notified := false
	bus.Subscribe("test-sub", func(event *IAxiomaticEvent) {
		notified = true
	})

	bus.Publish(&IAxiomaticEvent{ID: "evt-1", Type: "TEST"})

	if !notified {
		t.Error("expected subscriber to be notified")
	}

	notified = false
	bus.Unsubscribe("test-sub")
	bus.Publish(&IAxiomaticEvent{ID: "evt-2", Type: "TEST"})

	if notified {
		t.Error("expected subscriber not to be notified after unsubscription")
	}
}

func TestBaalAalEngine_Rules(t *testing.T) {
	engine := NewBaalAalEngine()
	called := 0
	engine.RegisterRule("TEST_EVENT", func(event *IAxiomaticEvent) {
		called++
	})

	engine.EventBus.Publish(&IAxiomaticEvent{Type: "TEST_EVENT"})
	engine.EventBus.Publish(&IAxiomaticEvent{Type: "OTHER_EVENT"})
	engine.EventBus.Publish(&IAxiomaticEvent{Type: "TEST_EVENT"})

	engine.ProcessCycle(1)

	if called != 2 {
		t.Errorf("expected rule to be called 2 times, got %d", called)
	}

	engine.ProcessCycle(2)
	if called != 2 {
		t.Errorf("expected no additional calls on second cycle, got %d", called)
	}
}

func TestBaalAalEngine_RecursivePublish(t *testing.T) {
	engine := NewBaalAalEngine()
	engine.RegisterRule("TRIGGER", func(event *IAxiomaticEvent) {
		engine.EventBus.Publish(&IAxiomaticEvent{Type: "RESPONSE"})
	})

	called := 0
	engine.RegisterRule("RESPONSE", func(event *IAxiomaticEvent) {
		called++
	})

	engine.EventBus.Publish(&IAxiomaticEvent{Type: "TRIGGER"})

	// First cycle processes TRIGGER, which publishes RESPONSE
	engine.ProcessCycle(1)

	// Second cycle processes RESPONSE
	engine.ProcessCycle(2)

	if called != 1 {
		t.Errorf("expected recursive rule to be called, got %d", called)
	}
}

func TestKappaSystem_HandleMove(t *testing.T) {
	ks := NewKappaSystem(nil)
	event := &IAxiomaticEvent{
		Metadata: map[string]interface{}{
			"client_id": "player-1",
			"x":         10.5,
			"y":         20.7,
		},
	}

	ks.HandleMove(event)

	pos, ok := ks.Positions["player-1"]
	if !ok {
		t.Fatal("expected position for player-1")
	}

	if pos[0] != 10500 || pos[1] != 20700 {
		t.Errorf("expected [10500, 20700], got %v", pos)
	}
}
