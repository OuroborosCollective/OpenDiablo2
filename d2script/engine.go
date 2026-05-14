package d2script

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore" // This causes the runtime to support underscore.js
)

const (
	defaultEvalTimeout = 500 * time.Millisecond
)

var (
	ErrEvalTimeout = errors.New("execution timed out")
)

// ScriptEngine allows running JavaScript scripts and Axiomatic logic.
type ScriptEngine struct {
	vm            *otto.Otto
	isEvalAllowed bool

	// BaalAal Engine
	BaalAal *BaalAalEngine
}

// CreateScriptEngine creates the script engine and returns a pointer to it.
func CreateScriptEngine() *ScriptEngine {
	vm := otto.New()
	err := vm.Set("debugPrint", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Script: %s\n", call.Argument(0).String())
		return otto.Value{}
	})

	if err != nil {
		fmt.Printf("could not bind the 'debugPrint' to the given function in script engine")
	}

	return &ScriptEngine{
		vm:            vm,
		isEvalAllowed: false,
		BaalAal:      NewBaalAalEngine(),
	}
}

// AllowEval allows the evaluation of JS code.
func (s *ScriptEngine) AllowEval() {
	s.isEvalAllowed = true
}

// DisallowEval disallows the evaluation of JS code.
func (s *ScriptEngine) DisallowEval() {
	s.isEvalAllowed = false
}

// ToValue converts the given interface{} value to a otto.Value
func (s *ScriptEngine) ToValue(source interface{}) (otto.Value, error) {
	return s.vm.ToValue(source)
}

// AddFunction adds the given function to the script engine with the given name.
func (s *ScriptEngine) AddFunction(name string, value interface{}) {
	err := s.vm.Set(name, value)
	if err != nil {
		fmt.Printf("could not add the '%s' function to the script engine", name)
	}
}

// RunScript runs the script file within the given path.
func (s *ScriptEngine) RunScript(fileName string) (*otto.Value, error) {
	fileData, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		fmt.Printf("could not read script file: %s\n", err.Error())
		return nil, err
	}

	val, err := s.vm.Run(string(fileData))
	if err != nil {
		fmt.Printf("Error running script: %s\n", err.Error())
		return nil, err
	}

	return &val, nil
}

// Eval JS code with a timeout.
func (s *ScriptEngine) Eval(code string) (res string, err error) {
	if !s.isEvalAllowed {
		return "", errors.New("disabled")
	}

	interrupt := make(chan func(), 1)
	s.vm.Interrupt = interrupt

	go func() {
		time.Sleep(defaultEvalTimeout)
		interrupt <- func() {
			panic(ErrEvalTimeout)
		}
	}()

	defer func() {
		if caught := recover(); caught != nil {
			if caught == ErrEvalTimeout {
				err = ErrEvalTimeout
				return
			}
			panic(caught)
		}
	}()

	val, err := s.vm.Eval(code)
	if err != nil {
		return "", err
	}

	return val.String(), nil
}
