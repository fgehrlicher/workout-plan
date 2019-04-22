package models

type Iteration struct {
	Type    string `yaml:"type"`
	Percent string `yaml:"percent"`
	Sets    string `yaml:"sets"`
	Reps    string `yaml:"reps"`
}
