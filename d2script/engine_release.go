//go:build release

package d2script

import (
	"errors"

	"github.com/robertkrimen/otto"
)

func (s *ScriptEngine) initJS() {
	// No-op for release builds
}

// ToValue returns an error in release builds.
func (s *ScriptEngine) ToValue(source interface{}) (otto.Value, error) {
	return otto.Value{}, errors.New("scripting disabled in release")
}

// AddFunction does nothing in release builds.
func (s *ScriptEngine) AddFunction(name string, value interface{}) {
}

// RunScript returns an error in release builds.
func (s *ScriptEngine) RunScript(fileName string) (*otto.Value, error) {
	return nil, errors.New("scripting disabled in release")
}

// Eval returns an error in release builds.
func (s *ScriptEngine) Eval(code string) (string, error) {
	return "", errors.New("scripting disabled in release")
}
