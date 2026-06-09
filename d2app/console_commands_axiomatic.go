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
		{"are-status", "show ARE-Logik system status", nil, a.areStatus},
		{"are-emerge", "trigger manual emergence cycle", nil, a.areEmerge},
		{"are-chunks", "show KAPPA chunk registry", nil, a.areChunks},
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

// areStatus displays the current ARE-Logik system status
func (a *App) areStatus([]string) error {
	resonance, expansion, entropy, tick := a.scriptEngine.GetAREStatus()

	a.terminal.Infof("=== Ouroboros ARE-Logik Status ===")
	a.terminal.Infof("Tick: %d", tick)
	a.terminal.Infof("Global Resonance: %.6f", resonance)
	a.terminal.Infof("Expansion Factor: %.6f", expansion)
	a.terminal.Infof("Entropy: %.6f", entropy)

	// Get BaalAal status
	baalRes, baalCycle := a.scriptEngine.BaalAal.GetStatus()
	a.terminal.Infof("BaalAal Resonance: %.6f", baalRes)
	a.terminal.Infof("BaalAal Cycle: %.6f", baalCycle)

	return nil
}

// areEmerge triggers a manual emergence cycle
func (a *App) areEmerge([]string) error {
	if a.ouroboros != nil {
		a.ouroboros.Advance()
	}
	resonance, expansion, entropy, tick := a.scriptEngine.GetAREStatus()

	a.terminal.Infof("Manual emergence cycle triggered")
	a.terminal.Infof("Tick: %d, Resonance: %.6f, Expansion: %.6f, Entropy: %.6f",
		tick, resonance, expansion, entropy)

	return nil
}

// areChunks displays the KAPPA chunk registry
func (a *App) areChunks([]string) error {
	if a.ouroboros == nil {
		a.terminal.Infof("Ouroboros system not initialized.")
		return nil
	}

	ouroboros := a.ouroboros

	ouroboros.mu.RLock()
	defer ouroboros.mu.RUnlock()

	chunkCount := len(ouroboros.chunkRegistry)
	a.terminal.Infof("=== KAPPA Chunk Registry (%d chunks) ===", chunkCount)

	count := 0
	for id, chunk := range ouroboros.chunkRegistry {
		chunk.mu.RLock()
		a.terminal.Infof("  %s: KAPPA(%d, %d), Resonance=%.3f, Occupants=%d, Gen=%d",
			id, chunk.X, chunk.Y, chunk.Resonance, len(chunk.Occupants), chunk.Generation)
		chunk.mu.RUnlock()

		count++
		if count >= 20 {
			a.terminal.Infof("  ... and %d more chunks", chunkCount-count)
			break
		}
	}

	return nil
}
