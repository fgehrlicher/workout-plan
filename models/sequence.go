package models

type Sequence struct {
	Type    string `yaml:"type" json:"type"`
	Percent string `yaml:"percent" json:"percent"`
	Sets    string `yaml:"sets" json:"sets"`
	Reps    string `yaml:"reps" json:"reps"`
}
