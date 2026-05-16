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
	vm            interface{}
	isEvalAllowed bool

	// BaalAal Engine
	BaalAal *BaalAalEngine
}

// CreateScriptEngine creates the script engine and returns a pointer to it.
func CreateScriptEngine() *ScriptEngine {
	s := &ScriptEngine{
		isEvalAllowed: false,
		BaalAal:      NewBaalAalEngine(),
	}
	s.initJS()
	return s
}

// AllowEval allows the evaluation of JS code.
func (s *ScriptEngine) AllowEval() {
	s.isEvalAllowed = true
}

// DisallowEval disallows the evaluation of JS code.
func (s *ScriptEngine) DisallowEval() {
	s.isEvalAllowed = false
}

// DispatchEvent dispatches an axiomatic event to the BaalAal engine.
func (s *ScriptEngine) DispatchEvent(event *IAxiomaticEvent) {
	s.BaalAal.EventBus.Publish(event)
}
