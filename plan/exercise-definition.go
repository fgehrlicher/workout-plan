package plan

import (
	"fmt"
)

type ExerciseDefinition struct {
	Id         string `yaml:"id" json:"id"`
	Name       string `yaml:"name" json:"name"`
	Definition struct {
		Video struct {
			Id string `yaml:"id" json:"id"`
		} `yaml:"video" json:"video"`
		Description string
	} `yaml:"definition" json:"definition"`
}

func (exerciseDefinition *ExerciseDefinition) Validate() error {
	if exerciseDefinition.Id == "" {
		return fmt.Errorf(
			"Id is required for exercise definitions. \nFull element: %+v",
			exerciseDefinition,
		)
	}

	return nil
}
