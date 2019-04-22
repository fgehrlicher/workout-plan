package models

type Unit struct {
	Name      string `yaml:"name" json:"name"`
	Exercises []Exercise `yaml:"exercises" json:"exercises"`
}