package models

type Exercise struct {
	Type     string     `yaml:"type" json:"type"`
	Exercise string     `yaml:"exercise" json:"exercise"`
	Sequence []Sequence `yaml:"sequence" json:"sequence"`
}
