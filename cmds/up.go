package cmds

import (
	"context"
	"fmt"
	"os"

	"github.com/environment-toolkit/et/internal/app"
	"github.com/environment-toolkit/et/internal/config"
)

type UpCmd struct {
	File   string `short:"f" help:"File containing the spec" default:"spec.yml"`
	Env    string `name:"env" help:"Environment to deploy the spec to" required:"true"`
	Region string `name:"region" help:"Region to deploy the spec to" required:"true"`
}

func (cmd *UpCmd) Run(globals *Globals) error {
	ctx := context.Background()

	line, err := config.New(ctx, cmd)
	if err != nil {
		line.Log.Error("error creating config", err)
		return fmt.Errorf("error creating config: %v", err)
	}

	manager, err := app.NewManager(ctx, line.App)
	if err != nil {
		line.Log.Error("error creating manager", err)
		return fmt.Errorf("error creating manager: %v", err)
	}

	data, err := os.ReadFile(cmd.File)
	if err != nil {
		line.Log.Error("error reading file", err)
		return fmt.Errorf("error reading file: %v", err)
	}
	variables := map[string]string{}

	result, err := manager.Up(ctx, string(data), variables)
	if err != nil {
		line.Log.Error("error uping", err)
		return fmt.Errorf("error uping: %v", err)
	}

	if result == nil {
		line.Log.Error("error uping no result")
		return fmt.Errorf("error uping no result")
	}

	return nil
}
