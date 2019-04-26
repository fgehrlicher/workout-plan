package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"workout-plan/exercise-definitions"
)

type ExerciseValidator func(*Exercise) error

var possibleExerciseTypes = map[string]ExerciseValidator{
	"main-exercise":       nil,
	"special-exercise":    nil,
	"additional-exercise": nil,
}

var exerciseDefinitions *exercise_definitions.ExerciseDefinitions

type Exercise struct {
	ExerciseDefinition    *ExerciseDefinition
	Type                  string              `yaml:"type" json:"type"`
	RawExerciseDefinition string              `yaml:"exercise-definition" json:"exercise-definition"`
	Sequence              []ExerciseIteration `yaml:"sequence" json:"sequence"`
}

func (exercise *Exercise) UnmarshalJSON(data []byte) error {
	var (
		jsonData                 map[string]string
		err                      error
		exerciseDefinitionString = "exercise-definition"
	)

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, exercise)
	if err != nil {
		return err
	}

	exerciseDefinitionValue, isSet := jsonData[exerciseDefinitionString]
	if !isSet {
		return errors.New(
			fmt.Sprintf(
				"`%v` is not set, but required for Exercise definiton elements",
				exerciseDefinitionString,
			),
		)
	}

	if exerciseDefinitions == nil {
		exerciseDefinitions = exercise_definitions.GetInstance()
	}

	exerciseDefinition, err := exerciseDefinitions.Get(exerciseDefinitionValue)
	if err != nil {
		return err
	}

	exercise.ExerciseDefinition = exerciseDefinition
	return nil
}

func (exercise *Exercise) Validate() error {
	err := TypeNotEmptyValidator(*exercise)
	if err != nil {
		return err
	}

	if exercise.RawExerciseDefinition == "" {
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

	return TypeNotAllowedError(*exercise)
}
