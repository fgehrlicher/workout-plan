package models

import (
	"errors"
	"fmt"
)

type ExerciseIterationValidator func(*ExerciseIteration) error

var possibleExerciseIterationsTypes = map[string]ExerciseIterationValidator{
	"sets-reps":              SetsRepsExerciseIterationValidator,
	"sets-reps-weight-range": SetsRepsWeightRangeExerciseIterationValidator,
}

type ExerciseIteration struct {
	Type      string `yaml:"type" json:"type"`
	MinWeight string `yaml:"min-weight" json:"min-weight"`
	MaxWeight string `yaml:"max-weight" json:"max-weight"`
	Percent   string `yaml:"percent" json:"percent"`
	Sets      string `yaml:"sets" json:"sets"`
	Reps      string `yaml:"reps" json:"reps"`
}

func (exerciseIteration *ExerciseIteration) Validate() error {
	err := TypeNotEmptyValidator(exerciseIteration)
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

	return TypeNotAllowedError(exerciseIteration)
}

func SetsRepsExerciseIterationValidator(exerciseIteration *ExerciseIteration) error {
	if exerciseIteration.Reps == "" || exerciseIteration.Sets == "" {
		if exerciseIteration.Type == "" {
			return errors.New(
				fmt.Sprintf(
					"reps and sets must be set if the type is set to sets-reps.\nFull element: %+v",
					exerciseIteration,
				),
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
			return errors.New(
				fmt.Sprintf(
					"reps, sets, min-weight and max-weight must be set if the " +
						"type is set to sets-reps-weight-range.\nFull element: %+v",
					exerciseIteration,
				),
			)
		}
	}

	return nil
}
