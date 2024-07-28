package main

import (
	"et/cmds"

	"github.com/alecthomas/kong"
)

func main() {
	cli := cmds.Cli{
		Globals: cmds.Globals{},
	}

	ctx := kong.Parse(&cli,
		kong.Name("et"),
		kong.Description("A tool for managing environments"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "0.0.1",
		})
	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
