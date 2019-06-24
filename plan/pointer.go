package plan

import (
	"time"
)

func CreatePointer(plan *Plan, userId string) Pointer {
	planPointer := Pointer{
		PlanId:      plan.ID,
		PlanVersion: plan.Version,
		UserId:      userId,
		Started:     time.Now(),
	}

	planPointer.Position.Unit = 1
	planPointer.Position.Exercise = 1

	return planPointer
}

type Pointer struct {
	PlanId      string    `bson:"plan_id"`
	PlanVersion string    `bson:"plan_version"`
	UserId      string    `bson:"user_id"`
	Started     time.Time `bson:"started"`
	Moved       time.Time `bson:"moved"`
	Position    struct {
		Unit     int `bson:"unit"`
		Exercise int `bson:"exercise"`
	} `bson:"position"`
	Data map[string]int `bson:"data"`
}
