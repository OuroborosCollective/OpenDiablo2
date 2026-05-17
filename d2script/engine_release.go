//go:build release

package d2script

import (
	"errors"

	"github.com/robertkrimen/otto"
)

var (
	ErrScriptingDisabled = errors.New("arbitrary scripting disabled in release builds")
)

func (s *ScriptEngine) initJS() {
	// No-op for release builds
}

// ToValue returns an error in release builds.
func (s *ScriptEngine) ToValue(source interface{}) (otto.Value, error) {
	return otto.Value{}, ErrScriptingDisabled
}

// AddFunction does nothing in release builds.
func (s *ScriptEngine) AddFunction(name string, value interface{}) {
}

// RunScript returns an error in release builds.
func (s *ScriptEngine) RunScript(fileName string) (*otto.Value, error) {
	return nil, ErrScriptingDisabled
}

// Eval returns an error in release builds.
func (s *ScriptEngine) Eval(code string) (string, error) {
	return "", ErrScriptingDisabled
}
