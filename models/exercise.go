package models

type ExerciseValidator func(*Exercise) error

var possibleExerciseTypes = map[string]ExerciseValidator{
	"main-exercise": nil,
}

type Exercise struct {
	Type     string              `yaml:"type" json:"type"`
	Exercise string              `yaml:"exercise" json:"exercise"`
	Sequence []ExerciseIteration `yaml:"sequence" json:"sequence"`
}

func (exercise *Exercise) Validate() error {
	err := TypeNotEmptyValidator(exercise)
	if err != nil {
		return err
	}

	for _, exerciseIteration := range exercise.Sequence {
		err = exerciseIteration.Validate()
		if err != nil {
			return err
		}
	}

	for possibleExerciseTypes, validator := range possibleExerciseTypes {
		if possibleExerciseTypes == exercise.Type {
			if validator != nil {
				return validator(exercise)
			}
		}
	}

	return TypeNotAllowedError(exercise)
}
