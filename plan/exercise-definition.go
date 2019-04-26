package plan

import (
	"errors"
	"fmt"
)

type ExerciseDefinition struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
}

func (exerciseDefinition *ExerciseDefinition) Validate() error {
	if exerciseDefinition.Name == "" {
		return errors.New(
			fmt.Sprintf(
				"Name is required for exercise definitions. \nFull element: %+v",
				exerciseDefinition,
			),
		)
	}

	return nil
}
