package plan

import (
	"fmt"
)

type ExerciseDefinition struct {
	Id          string `yaml:"id" json:"id"`
	Description string `yaml:"description" json:"description"`
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
