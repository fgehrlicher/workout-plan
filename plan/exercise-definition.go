package plan

import (
	"fmt"
)

type ExerciseDefinition struct {
	Id          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Video       struct {
		Id string `yaml:"id" json:"id"`
	} `yaml:"video" json:"video"`
	Image struct {
		Url string `yaml:"url" json:"url"`
	} `yaml:"video" json:"video"`
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
