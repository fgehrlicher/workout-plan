package models

type Plan struct {
	ID      string `yaml:"id" json:"id"`
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
	Units   []Unit `yaml:"units" json:"units"`
}
