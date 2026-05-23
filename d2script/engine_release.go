//go:build release

package d2script

import (
	"errors"
)

func (s *ScriptEngine) initJS() {
	// No-op for release builds
}

// ToValue returns an error in release builds.
func (s *ScriptEngine) ToValue(source interface{}) (interface{}, error) {
	return nil, errors.New("scripting disabled in release")
}

// AddFunction registers a handler in the BaalAal engine for release builds.
func (s *ScriptEngine) AddFunction(name string, value interface{}) {
	s.BaalAal.RegisterHandler(name, value)
}

// RunScript returns an error in release builds.
func (s *ScriptEngine) RunScript(fileName string) (*otto.Value, error) {
	return nil, errors.New("scripting disabled in release")
}

// Eval returns an error in release builds.
func (s *ScriptEngine) Eval(code string) (string, error) {
	return "", errors.New("scripting disabled in release")
}
