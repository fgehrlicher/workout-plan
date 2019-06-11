package plan

type ExerciseIteration struct {
	Type      string `yaml:"type" json:"type"`
	MinWeight string `yaml:"min-weight" json:"min-weight,omitempty"`
	MaxWeight string `yaml:"max-weight" json:"max-weight,omitempty"`
	Weight    string `yaml:"weight" json:"weight,omitempty"`
	Sets      string `yaml:"sets" json:"sets,omitempty"`
	Reps      string `yaml:"reps" json:"reps,omitempty"`
	Time      string `yaml:"time" json:"time,omitempty"`
	EmomTime  string `yaml:"emom-time" json:"emom-time,omitempty"`
	Variable  string `yaml:"variable" json:"variable,omitempty"`
}

func (exerciseIteration *ExerciseIteration) Validate() error {
	return TypeNotEmptyValidator(*exerciseIteration)
}
