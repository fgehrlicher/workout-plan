package models

type ExerciseValidator func(*Exercise) error

var possibleExerciseTypes = map[string]ExerciseValidator{
	"main-exercise": MainExerciseExerciseValidator,
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
			return validator(exercise)
		}
	}

	return TypeNotAllowedError(exercise)
}

func MainExerciseExerciseValidator(*Exercise) error {
	// No main exercise definition has been defined yet
	return nil
}
