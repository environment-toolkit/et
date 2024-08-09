package spec

type Props []byte

type Overrides struct {
	Environment string `yaml:"environment,omitempty"`
	Region      string `yaml:"region,omitempty"`
	Props       Props  `yaml:"props"`
}

type Resource struct {
	Type      string      `yaml:"type"`
	Name      string      `yaml:"name"`
	Props     Props       `yaml:"props"`
	Overrides []Overrides `yaml:"overrides,omitempty"`
}

type Access struct {
	Inbound  []string `yaml:"inbound,omitempty"`
	Outbound []string `yaml:"outbound,omitempty"`
}

type Spec struct {
	Resource

	Resources []Resource `yaml:"resources,omitempty"`
	Access    *Access    `yaml:"access,omitempty"`
}
