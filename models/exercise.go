package models

type Exercise struct {
	Type       string `yaml:"type"`
	Exercise   string `yaml:"exercise"`
	Iterations []Iteration `yaml:"iterations"`
}
