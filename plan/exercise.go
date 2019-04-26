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
	ExerciseDefinition *ExerciseDefinition
	Type               string
	Sequence           []ExerciseIteration
}

func (exercise *Exercise) UnmarshalJSON(data []byte) error {
	var (
		jsonData                       map[string]*json.RawMessage
		err                            error
		exerciseDefinitionResultString string
		exerciseDefinitionString       = "exercise-definition"
		sequenceString                 = "sequence"
		typeString                     = "type"
	)

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	typeValue, isSet := jsonData[typeString]
	if !isSet {
		return errors.New(
			fmt.Sprintf(
				"`%v` is not set, but required for Exercise elements",
				typeString,
			),
		)
	}

	err = json.Unmarshal(*typeValue, &exercise.Type)
	if err != nil {
		return err
	}

	sequenceValue, isSet := jsonData[sequenceString]
	if !isSet {
		return errors.New(
			fmt.Sprintf(
				"`%v` is not set, but required for Exercise elements",
				sequenceString,
			),
		)
	}

	err = json.Unmarshal(*sequenceValue, &exercise.Sequence)
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

	err = json.Unmarshal(*exerciseDefinitionValue, &exerciseDefinitionResultString)
	if err != nil {
		return err
	}

	exerciseDefinition, err := exerciseDefinitionsSingleton.Get(exerciseDefinitionResultString)
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
