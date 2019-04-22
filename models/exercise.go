package models

type Exercise struct {
	Type       string `yaml:"type" json:"type"`
	Exercise   string `yaml:"exercise" json:"exercise"`
	Iterations []Iteration `yaml:"iterations" json:"iterations"`
}
