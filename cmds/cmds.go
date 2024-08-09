package cmds

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type Globals struct {
	LogLevel string      `short:"l" help:"Set the logging level (debug|info|warn|error|fatal)" default:"info"`
	Version  VersionFlag `name:"version" help:"Print version information and quit"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type Cli struct {
	Globals

	Up   UpCmd   `cmd:"" help:"Create infrastructure from spec"`
	Down DownCmd `cmd:"" help:"Destroy infrastructure from spec"`
}
