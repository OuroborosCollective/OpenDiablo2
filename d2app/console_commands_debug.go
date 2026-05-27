//go:build !release
// +build !release

package d2app

func (a *App) initDebugCommands() {
	debugCommands := []struct {
		name string
		desc string
		args []string
		fn   func(args []string) error
	}{
		// Add debug commands here
	}

	for _, cmd := range debugCommands {
		if err := a.terminal.Bind(cmd.name, cmd.desc, cmd.args, cmd.fn); err != nil {
			a.Fatalf("failed to bind debug action %q: %v", cmd.name, err.Error())
		}
	}
}
