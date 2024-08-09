package cmds

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/environment-toolkit/et/pkg/spec"
)

type UpCmd struct {
	File   string `short:"f" help:"File containing the spec" default:"spec.yml"`
	Env    string `name:"env" help:"Environment to deploy the spec to" required:"true"`
	Region string `name:"region" help:"Region to deploy the spec to" required:"true"`
}

func (cmd *UpCmd) Run(globals *Globals) error {
	f, err := os.Open(cmd.File)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	defer f.Close()

	spec := &spec.Spec{}
	if err := yaml.NewDecoder(f).Decode(spec); err != nil {
		return fmt.Errorf("error unmarshalling yaml: %v", err)
	}

	// target := &models.Target{
	// 	Environment: cmd.Env,
	// 	Region:      cmd.Region,
	// }

	// // resolver
	// resolverManager := resolver.NewManager()

	// // parse the spec into a model. (this is to add all the overrides)
	// m, err := parser.Parse(resolverManager, target, spec)
	// if err != nil {
	// 	return fmt.Errorf("error parsing spec: %v", err)
	// }

	// find all ${{ .. }} values

	// download the spec file

	// resolve all values
	fmt.Printf("Resolved spec: %v\n", spec)

	return nil
}
