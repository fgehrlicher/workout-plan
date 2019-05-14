package plan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ExerciseDefinitions struct {
	underlyingSlice []ExerciseDefinition
}

func (exerciseDefinitions *ExerciseDefinitions) Add(exerciseDefinition ExerciseDefinition) {
	logEntry := log.WithFields(log.Fields{
		"Id": exerciseDefinition.Name,
	})

	for _, existingExerciseDefinition := range exerciseDefinitions.underlyingSlice {
		if existingExerciseDefinition.Name == exerciseDefinition.Name {
			logEntry.Warning("Exercise definition already exists")
			return
		}
	}

	exerciseDefinitions.underlyingSlice = append(exerciseDefinitions.underlyingSlice, exerciseDefinition)
	logEntry.Info("Exercise definition added")
}

func (exerciseDefinitions *ExerciseDefinitions) Get(name string) (*ExerciseDefinition, error) {
	for _, exerciseDefinition := range exerciseDefinitions.underlyingSlice {
		if exerciseDefinition.Name == name {
			return &exerciseDefinition, nil
		}
	}

	return nil, errors.New(
		fmt.Sprintf(
			"exercise definition with name `%v` was not found.",
			name,
		),
	)
}

func (exerciseDefinitions *ExerciseDefinitions) GetAll() []*ExerciseDefinition {
	var returnDefinitions []*ExerciseDefinition
	for key := range exerciseDefinitions.underlyingSlice {
		returnDefinitions = append(returnDefinitions, &exerciseDefinitions.underlyingSlice[key])
	}

	return returnDefinitions
}

var exerciseDefinitionsSingleton *ExerciseDefinitions
var exerciseDefinitionsOnce sync.Once

func GetExerciseDefinitionsInstance() *ExerciseDefinitions {
	exerciseDefinitionsOnce.Do(func() {
		exerciseDefinitionsSingleton = &ExerciseDefinitions{}
	})
	return exerciseDefinitionsSingleton
}

func InitializeExerciseDefinitions(exerciseDefinitionFile string) error {
	exerciseDefinitions := GetExerciseDefinitionsInstance()

	fileExtension := filepath.Ext(exerciseDefinitionFile)
	fileData, err := ioutil.ReadFile(exerciseDefinitionFile)
	if err != nil {
		return errors.New(
			fmt.Sprintf(
				"can´t load exercise definition file (tried: '%v'): %v",
				exerciseDefinitionFile,
				err,
			),
		)
	}

	var exerciseDefinitionsSlice []ExerciseDefinition

	switch fileExtension {
	case ".yml":
		err = yaml.Unmarshal(fileData, &exerciseDefinitionsSlice)
		if err != nil {
			return err
		}
	case ".json":
		err = json.Unmarshal(fileData, &exerciseDefinitionsSlice)
		if err != nil {
			return err
		}
	default:
		return errors.New(
			fmt.Sprintf(
				"Invalid file extension: `%v",
				fileExtension,
			),
		)
	}

	for _, exerciseDefinitionsElement := range exerciseDefinitionsSlice {
		exerciseDefinitions.Add(exerciseDefinitionsElement)
	}

	return nil
}
