package plan

import (
	"errors"
	"fmt"
)

type Unit struct {
	Name      string     `yaml:"name" json:"name"`
	Exercises []Exercise `yaml:"exercises" json:"exercises"`
}

func (unit *Unit) Validate() error {
	if unit.Name == "" {
		return errors.New(
			fmt.Sprintf(
				"the name field of units musn´t be empty.\nFull element: %+v",
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

func (unit *Unit) GetRequiredVariables() []string {
	var variables []string

	for _, exercise := range unit.Exercises {
		for _, iteration := range exercise.Sequence {
			if iteration.Variable != "" {
				variables = append(variables, iteration.Variable)
			}
		}
	}

	return variables
}
