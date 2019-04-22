package models

type Unit struct {
	Name      string `yaml:"name"`
	Exercises []Exercise `yaml:"exercises"`
}