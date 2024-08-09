package models

type Target struct {
	Environment string
	Region      string
}

type Resource struct {
	Type  string
	Name  string
	Props map[string]interface{}
}

type Access struct {
	Inbound  []string
	Outbound []string
}

type Manifest struct {
	Resource

	Resources []*Resource
	Access    *Access
}
