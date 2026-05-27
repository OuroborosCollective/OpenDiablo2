package d2script

import (
	"errors"
	"time"
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
}

// CreateScriptEngine creates the script engine and returns a pointer to it.
func CreateScriptEngine() *ScriptEngine {
	s := &ScriptEngine{
		BaalAal: NewBaalAalEngine(),
	}
	s.initJS()
	return s
}

// DispatchEvent dispatches an axiomatic event to the BaalAal engine.
// This is the primary way logic should be handled in OpenDiablo2.
func (s *ScriptEngine) DispatchEvent(event *IAxiomaticEvent) {
	if s.BaalAal != nil {
		s.BaalAal.EventBus.Publish(event)
	}
}
