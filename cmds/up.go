package cmds

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/environment-toolkit/et/internal/app"
	"github.com/environment-toolkit/et/internal/config"
)

type UpCmd struct {
	File   string `short:"f" help:"File containing the spec" default:"spec.yml"`
	DryRun bool   `name:"dry-run" help:"Dry run the spec"`
	Env    string `name:"env" help:"Environment to deploy the spec to" required:"true"`
	Region string `name:"region" help:"Region to deploy the spec to" required:"true"`
}

func (cmd *UpCmd) Run(globals *Globals) error {
	ctx := context.Background()

	cfgLine, err := config.New(ctx, cmd)
	if err != nil {
		cfgLine.Log.Error("error creating config", err)
		return fmt.Errorf("error creating config: %v", err)
	}

	data, err := os.ReadFile(cmd.File)
	if err != nil {
		cfgLine.Log.Error("error reading file", err)
		return fmt.Errorf("error reading file: %v", err)
	}
	filename := filepath.Base(cmd.File)
	variables := map[string]string{}

	appManager, err := app.NewManager(ctx, cfgLine.App)
	if err != nil {
		cfgLine.Log.Error("error creating app manager", err)
		return fmt.Errorf("error creating app manager: %v", err)
	}

	upConfig := app.UpConfig{
		DryRun:    cmd.DryRun,
		Spec:      data,
		Filename:  filename,
		Variables: variables,
	}

	result, err := appManager.Up(ctx, upConfig)
	if err != nil {
		cfgLine.Log.Error("error uping", err)
		return fmt.Errorf("error uping: %v", err)
	}

	if result == nil {
		cfgLine.Log.Error("error uping no result")
		return fmt.Errorf("error uping no result")
	}

	return nil
}
