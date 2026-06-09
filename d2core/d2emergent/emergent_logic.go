package d2emergent

import (
	"fmt"
	"math"
	"time"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

// ARELogikEngine represents the axiomatic Rekursion framework for emergent NPC behavior.
type ARELogikEngine struct {
	logger          *d2util.Logger
	GlobalResonance float64
	Tick            uint64
	Expansion       float64
	Entropy         float64
	Feedback        float64
}

// Markgraf constants for ARE-Logik (re-exported for backwards compatibility)
const (
	MarkgrafHarmonicFrequency = 0.017
	MarkgrafExpansionCoeff    = 1.0000001
	MarkgrafEntropyDecay      = 0.999
	MarkgrafFeedbackWeight    = 0.05
)

func CreateARELogikEngine(l d2util.LogLevel) *ARELogikEngine {
	logger := d2util.NewLogger()
	logger.SetPrefix("ARE-Logik")
	logger.SetLevel(l)
	return &ARELogikEngine{
		logger:          logger,
		GlobalResonance: 1.0,
		Tick:            0,
		Expansion:       1.0,
		Entropy:         0.5,
		Feedback:        0.0,
	}
}

func (e *ARELogikEngine) ProcessEmergence() *d2script.IAxiomaticEvent {
	// Implement Ouroboros Collective Markgraf ARE-Logik
	// This handles the endless world expansion and living NPC logic via resonance fluctuation.

	e.Tick++

	// Markgraf Harmonic Oscillation: A combination of multiple sine waves for deterministic complexity
	harmonic := math.Sin(float64(e.Tick)*MarkgrafHarmonicFrequency) +
		0.5*math.Cos(float64(e.Tick)*MarkgrafHarmonicFrequency*0.618)

	// Ouroboros Feedback: The previous resonance influences the current state
	e.GlobalResonance = math.Abs(harmonic + (e.Feedback * MarkgrafFeedbackWeight))

	// Update feedback loop
	e.Feedback = e.GlobalResonance

	// Endless World Expansion logic
	e.Expansion *= MarkgrafExpansionCoeff
	if e.Expansion > 1e9 { // Safety reset for extreme expansion
		e.Expansion = 1.0
	}

	// Living NPC Entropy (Liveliness)
	e.Entropy = (e.Entropy * MarkgrafEntropyDecay) + (e.GlobalResonance * (1.0 - MarkgrafEntropyDecay))

	if e.logger != nil {
		e.logger.Debugf("ARE-Logik: Tick=%d, Resonance=%.6f, Expansion=%.6f, Entropy=%.6f",
			e.Tick, e.GlobalResonance, e.Expansion, e.Entropy)
	}

	return &d2script.IAxiomaticEvent{
		ID:        fmt.Sprintf("WorldEmergence-%d", time.Now().UnixNano()),
		Type:      "WorldEmergence",
		Timestamp: time.Now().Unix(),
		Payload:   e.GlobalResonance,
		Metadata: map[string]interface{}{
			"expansion": e.Expansion,
			"entropy":   e.Entropy,
			"tick":      e.Tick,
			"logik":     "Markgraf",
		},
	}
}
