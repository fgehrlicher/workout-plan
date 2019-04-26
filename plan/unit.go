package plan

import (
	"errors"
	"fmt"
)

type Unit struct {
	Name      string `yaml:"name" json:"name"`
	Exercises []Exercise `yaml:"exercises" json:"exercises"`
}

func (unit *Unit) Validate() error {
	if unit.Name == "" {
		return errors.New(
			fmt.Sprintf(
				"the name of units musn´t be empty.\nFull element: %+v",
				unit,
			),
		)
	}

	if len(unit.Exercises) == 0 {
		return errors.New(
			fmt.Sprintf(
				"the unit with name %v does not have any exercises",
				unit.Name,
			),
		)
	}

	for _, exercise := range unit.Exercises {
		err := exercise.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
