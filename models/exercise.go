package models

type Exercise struct {
	Type     string              `yaml:"type" json:"type"`
	Exercise string              `yaml:"exercise" json:"exercise"`
	Sequence []ExerciseIteration `yaml:"sequence" json:"sequence"`
}
