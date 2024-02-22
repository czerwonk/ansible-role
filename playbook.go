package main

type playbook struct {
	Hosts       []string `yaml:"hosts"`
	GatherFacts bool     `yaml:"gather_facts"`
	Roles       []string `yaml:"roles"`
}
