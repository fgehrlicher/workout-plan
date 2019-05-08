package plan

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
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
		return FieldRequiredForExerciseError(typeString)
	}

	err = json.Unmarshal(*typeValue, &exercise.Type)
	if err != nil {
		return err
	}

	sequenceValue, isSet := jsonData[sequenceString]
	if !isSet {
		return FieldRequiredForExerciseError(sequenceString)
	}

	err = json.Unmarshal(*sequenceValue, &exercise.Sequence)
	if err != nil {
		return err
	}

	exerciseDefinitionValue, isSet := jsonData[exerciseDefinitionString]
	if !isSet {
		return FieldRequiredForExerciseError(exerciseDefinitionString)
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

func (exercise *Exercise) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var (
		yamlData                       map[string]interface{}
		err                            error
		exerciseDefinitionResultString string
		exerciseDefinitionString       = "exercise-definition"

		// That temp struct is needed because there is no way (to my knowledge)
		// to unmarshal the whole Exercise to something like a slice of
		// json RawMessages and then unmarshal the Sequence further.
		// @TODO replace that hack
		tempStruct = struct {
			Type     string              `yaml:"type"`
			Sequence []ExerciseIteration `yaml:"sequence"`
		}{}
	)

	err = unmarshal(&tempStruct)
	if err != nil {
		return err
	}

	exercise.Sequence = tempStruct.Sequence
	exercise.Type = tempStruct.Type

	err = unmarshal(&yamlData)
	if err != nil {
		return err
	}

	exerciseDefinitionValue, isSet := yamlData[exerciseDefinitionString]
	if !isSet {
		return FieldRequiredForExerciseError(exerciseDefinitionString)
	}

	err = yaml.Unmarshal([]byte(exerciseDefinitionValue.(string)), &exerciseDefinitionResultString)
	if err != nil {
		return err
	}

	if exerciseDefinitionsSingleton == nil {
		GetExerciseDefinitionsInstance()
	}

	exerciseDefinition, err := exerciseDefinitionsSingleton.Get(exerciseDefinitionResultString)
	if err != nil {
		return err
	}

	exercise.ExerciseDefinition = exerciseDefinition
	return nil
}

func FieldRequiredForExerciseError(fieldName string) error {
	return errors.New(
		fmt.Sprintf(
			"`%v` is not set, but required for Exercise elements",
			fieldName,
		),
	)
}

func (exercise *Exercise) Validate() error {
	err := TypeNotEmptyValidator(*exercise)
	if err != nil {
		return err
	}

	if exercise.ExerciseDefinition == nil {
		return FieldRequiredForExerciseError("definitions-file")
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
