package models

import (
	"errors"
	"fmt"
)

type ExerciseValidator func(*Exercise) error

var possibleExerciseTypes = map[string]ExerciseValidator{
	"main-exercise":       nil,
	"special-exercise":    nil,
	"additional-exercise": nil,
}

type Exercise struct {
	Type               string              `yaml:"type" json:"type"`
	ExerciseDefinition string              `yaml:"exercise-definition" json:"exercise-definition"`
	Sequence           []ExerciseIteration `yaml:"sequence" json:"sequence"`
}

func (exercise *Exercise) Validate() error {
	err := TypeNotEmptyValidator(exercise)
	if err != nil {
		return err
	}

	if exercise.ExerciseDefinition == "" {
		return errors.New(
			fmt.Sprintf(
				"the exercise definition musn´t be empty for exercise elements.\nFull element: %+v",
				exercise,
			),
		)
	}

	if len(exercise.Sequence) == 0 {
		return errors.New(
			fmt.Sprintf(
				"exercises must have an sequence with at least one item.\nFull element: %+v",
				exercise,
			),
		)
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
			return nil
		}
	}

	return TypeNotAllowedError(exercise)
}
