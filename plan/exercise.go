package plan

import (
	"encoding/json"
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

	if exerciseDefinitionsSingleton == nil {
		GetExerciseDefinitionsInstance()
	}

	exerciseDefinition, err := exerciseDefinitionsSingleton.Get(exerciseDefinitionValue)
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
