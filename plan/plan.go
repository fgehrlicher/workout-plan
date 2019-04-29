package plan

import (
	"errors"
	"fmt"
)

type Plan struct {
	ID      string `yaml:"id" json:"id"`
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
	Units   []Unit `yaml:"units" json:"units"`
}

func (plan *Plan) Validate() error {
	if plan.ID == "" || plan.Name == "" || plan.Version == "" {
		return errors.New(
			fmt.Sprintf(
				"the plan id, name and version must be set.\nFull element: %+v",
				plan,
			),
		)
	}

	if len(plan.Units) == 0 {
		return errors.New(
			fmt.Sprintf(
				"the plan with id %v does not have any units",
				plan.ID,
			),
		)
	}

	for _, unit := range plan.Units {
		err := unit.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}