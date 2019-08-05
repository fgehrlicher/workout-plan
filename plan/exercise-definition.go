package plan

import (
	"fmt"
)

type ExerciseDefinition struct {
	Id          string  `yaml:"id" json:"id"`
	Name        string  `yaml:"name" json:"name"`
	Description string  `yaml:"description" json:"description"`
	Media       []Media `yaml:"media" json:"media"`
}

type Media struct {
	Id  string `yaml:"id" json:"id,omitempty"`
	Url string `yaml:"url" json:"url,omitempty"`
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
