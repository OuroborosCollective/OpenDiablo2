package d2script

import (
	"errors"
	"time"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2emergent"
)

const (
	defaultEvalTimeout = 500 * time.Millisecond
)

var (
	ErrEvalTimeout = errors.New("execution timed out")
)

// ScriptEngine allows running JavaScript scripts and Axiomatic logic.
type ScriptEngine struct {
	vm interface{}

	// BaalAal Engine
	BaalAal *BaalAalEngine

	// Ouroboros ARE-Logik System - binds axiomatic logic to game systems
	Ouroboros *d2emergent.OuroborosLogikSystem

	// NPC Emergent System
	NPCEmergent *d2emergent.NPCEmergentSystem

	// Combat Emergent System
	CombatEmergent *d2emergent.CombatEmergentSystem

	// Item Emergent System
	ItemEmergent *d2emergent.ItemEmergentSystem

	// Log level for subsystem logging
	logLevel d2util.LogLevel
}

// CreateScriptEngine creates the script engine and returns a pointer to it.
func CreateScriptEngine() *ScriptEngine {
	return CreateScriptEngineWithLogLevel(d2util.LogLevelDefault)
}

// CreateScriptEngineWithLogLevel creates the script engine with a specific log level.
func CreateScriptEngineWithLogLevel(logLevel d2util.LogLevel) *ScriptEngine {
	s := &ScriptEngine{
		BaalAal:  NewBaalAalEngine(),
		logLevel: logLevel,
	}

	// Initialize Ouroboros-ARE-Logik system
	s.Ouroboros = d2emergent.NewOuroborosLogikSystem(s.BaalAal, logLevel)

	// Initialize NPC Emergent System
	s.NPCEmergent = d2emergent.NewNPCEmergentSystem(s.Ouroboros, s.BaalAal, logLevel)

	// Initialize Combat Emergent System
	s.CombatEmergent = d2emergent.NewCombatEmergentSystem(s.Ouroboros, s.BaalAal, logLevel)

	// Initialize Item Emergent System
	s.ItemEmergent = d2emergent.NewItemEmergentSystem(s.Ouroboros, s.BaalAal, logLevel)

	s.initJS()
	return s
}

// DispatchEvent dispatches an axiomatic event to the BaalAal engine.
// This is the primary way logic should be handled in OpenDiablo2.
func (s *ScriptEngine) DispatchEvent(event *IAxiomaticEvent) {
	if s.BaalAal != nil {
		s.BaalAal.EventBus.Publish(event)
	}

	// Also process through Ouroboros system
	if s.Ouroboros != nil {
		s.Ouroboros.ProcessEvent(event)
	}
}

// Advance advances all emergent systems
func (s *ScriptEngine) Advance() {
	if s.Ouroboros != nil {
		s.Ouroboros.Advance()
	}
	if s.NPCEmergent != nil {
		s.NPCEmergent.Advance()
	}
	if s.CombatEmergent != nil {
		s.CombatEmergent.Advance()
	}
	if s.ItemEmergent != nil {
		s.ItemEmergent.Advance()
	}
}

// GetAREStatus returns the current ARE-Logik system status
func (s *ScriptEngine) GetAREStatus() (resonance, expansion, entropy float64, tick uint64) {
	if s.Ouroboros != nil {
		return s.Ouroboros.GetStatus()
	}
	return 0, 0, 0, 0
}
