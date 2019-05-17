package plan

import (
	"fmt"
)

type ExerciseIterationValidator func(*ExerciseIteration) error

const (
	SetsRepsExerciseIterationType            = "sets-reps"
	SetsRepsWeightRangeExerciseIterationType = "sets-reps-weight-range"
	MaxOutExerciseIterationType              = "max-out"
)

var possibleExerciseIterationsTypes = map[string]ExerciseIterationValidator{
	SetsRepsExerciseIterationType:            SetsRepsExerciseIterationValidator,
	SetsRepsWeightRangeExerciseIterationType: SetsRepsWeightRangeExerciseIterationValidator,
	MaxOutExerciseIterationType:              nil,
}

type ExerciseIteration struct {
	Type      string `yaml:"type" json:"type"`
	MinWeight string `yaml:"min-weight" json:"min-weight,omitempty"`
	MaxWeight string `yaml:"max-weight" json:"max-weight,omitempty"`
	Weight    string `yaml:"weight" json:"weight,omitempty"`
	Sets      string `yaml:"sets" json:"sets,omitempty"`
	Reps      string `yaml:"reps" json:"reps,omitempty"`
	Variable  string `yaml:"variable" json:"variable,omitempty"`
}

func (exerciseIteration *ExerciseIteration) Validate() error {
	err := TypeNotEmptyValidator(*exerciseIteration)
	if err != nil {
		return err
	}

	for possibleExerciseIterationType, validator := range possibleExerciseIterationsTypes {
		if possibleExerciseIterationType == exerciseIteration.Type {
			if validator != nil {
				return validator(exerciseIteration)
			}
			return nil
		}
	}

	return TypeNotAllowedError(*exerciseIteration)
}

func SetsRepsExerciseIterationValidator(exerciseIteration *ExerciseIteration) error {
	if exerciseIteration.Reps == "" || exerciseIteration.Sets == "" {
		if exerciseIteration.Type == "" {
			return fmt.Errorf(
				"reps and sets must be set if the type is set to sets-reps.\nFull element: %+v",
				exerciseIteration,
			)
		}
	}

	return nil
}

func SetsRepsWeightRangeExerciseIterationValidator(exerciseIteration *ExerciseIteration) error {
	if exerciseIteration.Reps == "" ||
		exerciseIteration.Sets == "" ||
		exerciseIteration.MaxWeight == "" ||
		exerciseIteration.MinWeight == "" {
		if exerciseIteration.Type == "" {
			return fmt.Errorf(
				"reps, sets, min-weight and max-weight must be set if the "+
					"type is set to sets-reps-weight-range.\nFull element: %+v",
				exerciseIteration,
			)
		}
	}

	return nil
}
