package models

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

	for _, exercise := range unit.Exercises {
		err := exercise.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
