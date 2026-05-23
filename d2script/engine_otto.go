//go:build !release

package d2script

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

func (s *ScriptEngine) initJS() {
	vm := otto.New()
	err := vm.Set("debugPrint", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Script: %s\n", call.Argument(0).String())
		return otto.Value{}
	})

	if err != nil {
		fmt.Printf("could not bind the 'debugPrint' to the given function in script engine")
	}
	s.vm = vm
}

func (s *ScriptEngine) getVM() *otto.Otto {
	return s.vm.(*otto.Otto)
}

// ToValue converts the given interface{} value to a otto.Value
func (s *ScriptEngine) ToValue(source interface{}) (otto.Value, error) {
	return s.getVM().ToValue(source)
}

// AddFunction adds the given function to the script engine with the given name.
func (s *ScriptEngine) AddFunction(name string, value interface{}) {
	err := s.getVM().Set(name, value)
	if err != nil {
		fmt.Printf("could not add the '%s' function to the script engine", name)
	}
}

// RunScript runs the script file within the given path.
func (s *ScriptEngine) RunScript(fileName string) (val *otto.Value, err error) {
	fileData, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		fmt.Printf("could not read script file: %s\n", err.Error())
		return nil, err
	}

	vm := s.getVM()
	interrupt := make(chan func(), 1)
	vm.Interrupt = interrupt

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

	res, err := vm.Run(string(fileData))
	if err != nil {
		fmt.Printf("Error running script: %s\n", err.Error())
		return nil, err
	}

	return &res, nil
}

// Eval JS code with a timeout.
func (s *ScriptEngine) Eval(code string) (res string, err error) {
	if !s.isEvalAllowed {
		return "", errors.New("disabled")
	}

	vm := s.getVM()
	interrupt := make(chan func(), 1)
	vm.Interrupt = interrupt

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

	val, err := vm.Eval(code)
	if err != nil {
		return "", err
	}

	return val.String(), nil
}
