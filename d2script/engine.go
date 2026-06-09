package d2script

import (
	"errors"
	"time"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
)

const (
	defaultEvalTimeout = 500 * time.Millisecond
)

var (
	ErrEvalTimeout = errors.New("execution timed out")
)

// EmergentSystems holds references to emergent game systems.
// This interface allows the script engine to interact with emergent systems
// without creating an import cycle between d2script and d2emergent.
type EmergentSystems interface {
	Advance()
	GetAREStatus() (resonance, expansion, entropy float64, tick uint64)
	GetEntityState(entityID string) interface{}
}

// ScriptEngine allows running JavaScript scripts and Axiomatic logic.
type ScriptEngine struct {
	vm interface{}

	// BaalAal Engine
	BaalAal *BaalAalEngine

	// Emergent systems interface (set externally to avoid import cycle)
	Ouroboros EmergentSystems

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

	s.initJS()
	return s
}

// SetEmergentSystems sets the emergent systems interface after initialization.
// This breaks the import cycle between d2script and d2emergent.
func (s *ScriptEngine) SetEmergentSystems(emergent EmergentSystems) {
	s.Ouroboros = emergent
}

// DispatchEvent dispatches an axiomatic event to the BaalAal engine.
// This is the primary way logic should be handled in OpenDiablo2.
func (s *ScriptEngine) DispatchEvent(event *IAxiomaticEvent) {
	if s.BaalAal != nil {
		s.BaalAal.EventBus.Publish(event)
	}
}

// Advance advances all emergent systems
func (s *ScriptEngine) Advance() {
	if s.Ouroboros != nil {
		s.Ouroboros.Advance()
	}
}

// GetAREStatus returns the current ARE-Logik system status
func (s *ScriptEngine) GetAREStatus() (resonance, expansion, entropy float64, tick uint64) {
	if s.Ouroboros != nil {
		return s.Ouroboros.GetAREStatus()
	}
	return 0, 0, 0, 0
}
