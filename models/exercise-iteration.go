package models

import (
	"errors"
	"fmt"
)

type ExerciseIterationValidator func(*ExerciseIteration) error

var possibleExerciseIterationsTypes = map[string]ExerciseIterationValidator{
	"sets-reps":              ExerciseIterationSetsRepsValidator,
	"sets-reps-weight-range": ExerciseIterationSetsRepsWeightRangeValidator,
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

	for possibleExerciseIteration, validator := range possibleExerciseIterationsTypes {
		if possibleExerciseIteration == exerciseIteration.Type {
			return validator(exerciseIteration)
		}
	}

	return TypeNotAllowedError(exerciseIteration)
}

func ExerciseIterationSetsRepsValidator(exerciseIteration *ExerciseIteration) error {
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

func ExerciseIterationSetsRepsWeightRangeValidator(exerciseIteration *ExerciseIteration) error {
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
