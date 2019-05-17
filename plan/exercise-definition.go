package plan

import (
	"fmt"
)

type ExerciseDefinition struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
}

func (exerciseDefinition *ExerciseDefinition) Validate() error {
	if exerciseDefinition.Name == "" {
		return fmt.Errorf(
			"Id is required for exercise definitions. \nFull element: %+v",
			exerciseDefinition,
		)
	}

	return nil
}
