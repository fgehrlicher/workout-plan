package models

type Plan struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Units []Unit `yaml:"units"`
}
