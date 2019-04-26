package exercise_definitions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"workout-plan/models"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ExerciseDefinitions struct {
	underlyingSlice []models.ExerciseDefinition
}

func (exerciseDefinitions *ExerciseDefinitions) Add(exerciseDefinition models.ExerciseDefinition) {
	logEntry := log.WithFields(log.Fields{
		"Name": exerciseDefinition.Name,
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

func (exerciseDefinitions *ExerciseDefinitions) Get(name string) (*models.ExerciseDefinition, error) {
	for _, exerciseDefinition := range exerciseDefinitions.underlyingSlice {
		if exerciseDefinition.Name == name {
			return &exerciseDefinition, nil
		}
	}

	return nil, errors.New(
		fmt.Sprintf(
			"exercise definition with name`%v` was not found.",
			name,
		),
	)
}

var instance *ExerciseDefinitions
var once sync.Once

func GetInstance() *ExerciseDefinitions {
	once.Do(func() {
		instance = &ExerciseDefinitions{}
	})
	return instance
}

func InitializeExerciseDefinitions(exerciseDefinitionFile string) error {
	exerciseDefinitions := GetInstance()

	fileExtension := filepath.Ext(exerciseDefinitionFile)
	fileData, err := ioutil.ReadFile(exerciseDefinitionFile)
	if err != nil {
		return err
	}

	var exerciseDefinitionsSlice []models.ExerciseDefinition

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
