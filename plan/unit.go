package plan

import (
	"errors"
	"fmt"
)

type Unit struct {
	Id        string     `yaml:"id" json:"id"`
	Exercises []Exercise `yaml:"exercises" json:"exercises"`
}

func (unit *Unit) Validate() error {
	if unit.Id == "" {
		return errors.New(
			fmt.Sprintf(
				"the id field of units musn´t be empty.\nFull element: %+v",
				unit,
			),
		)
	}

	if len(unit.Exercises) == 0 {
		return errors.New(
			fmt.Sprintf(
				"the unit with id %v does not have any exercises",
				unit.Id,
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
