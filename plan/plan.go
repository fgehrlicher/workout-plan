package plan

import (
	"fmt"
)

type Plan struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Version     string `yaml:"version" json:"version"`
	Units       []Unit `yaml:"units" json:"units,omitempty"`
}

func (plan *Plan) Validate() error {
	if plan.ID == "" || plan.Name == "" || plan.Version == "" {
		return fmt.Errorf(
			"the plan id, name and version must be set.\nFull element: %+v",
			plan,
		)
	}

	if len(plan.Units) == 0 {
		return fmt.Errorf(
			"the plan with id %v does not have any units",
			plan.ID,
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

func (plan *Plan) GetSanitizedCopy() Plan {
	sanitizedPlan := *plan
	sanitizedPlan.Units = nil
	return sanitizedPlan
}
