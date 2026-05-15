// +build !release

package d2app

func (a *App) initDebugCommands() {
	debugCommands := []struct {
		name string
		desc string
		args []string
		fn   func(args []string) error
	}{
		{"js", "eval JS scripts", []string{"code"}, a.evalJS},
	}

	for _, cmd := range debugCommands {
		if err := a.terminal.Bind(cmd.name, cmd.desc, cmd.args, cmd.fn); err != nil {
			a.Fatalf("failed to bind debug action %q: %v", cmd.name, err.Error())
		}
	}
}

func (a *App) evalJS(args []string) error {
	val, err := a.scriptEngine.Eval(args[0])
	if err != nil {
		a.terminal.Errorf(err.Error())
		return nil
	}

	a.Info("%s" + val)

	return nil
}
