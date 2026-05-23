package d2app

import (
	"fmt"
	"time"

	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

func (a *App) initAxiomaticCommands() {
	axCommands := []struct {
		name string
		desc string
		args []string
		fn   func(args []string) error
	}{
		{"ax-pub", "publish an axiomatic event", []string{"type", "payload"}, a.axPublish},
		{"ax-history", "show axiomatic event history", nil, a.axHistory},
		{"ax-stats", "show axiomatic engine stats", nil, a.axStats},
	}

	for _, cmd := range axCommands {
		if err := a.terminal.Bind(cmd.name, cmd.desc, cmd.args, cmd.fn); err != nil {
			a.Fatalf("failed to bind axiomatic action %q: %v", cmd.name, err.Error())
		}
	}
}

func (a *App) axPublish(args []string) error {
	event := &d2script.IAxiomaticEvent{
		ID:        fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		Type:      args[0],
		Timestamp: time.Now().Unix(),
		Payload:   args[1],
	}

	a.scriptEngine.BaalAal.EventBus.Publish(event)
	a.terminal.Infof("Published axiomatic event: %s", event.ID)

	return nil
}

func (a *App) axHistory([]string) error {
	history := a.scriptEngine.BaalAal.EventBus.GetHistory()
	if len(history) == 0 {
		a.terminal.Infof("No history found.")
		return nil
	}

	for _, event := range history {
		if event == nil {
			continue
		}
		a.terminal.Infof("[%d] %s: %v",
			event.SequenceID, event.Type, event.Payload)
	}
	return nil
}

func (a *App) axStats([]string) error {
	a.terminal.Infof("Axiomatic Engine (BaalAal) is active.")
	return nil
}
