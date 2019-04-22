package models

type Plan struct {
	ID      string `yaml:"id"`
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Units   []Unit `yaml:"units"`
}
