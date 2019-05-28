package plan

func CreatePointer(plan *Plan, userId string) Pointer {
	planPointer := Pointer{
		PlanId:      plan.ID,
		PlanVersion: plan.Version,
		UserId:      userId,
	}

	planPointer.Position.Unit = 1
	planPointer.Position.Exercise = 1

	return planPointer
}

type Pointer struct {
	PlanId      string `bson:"plan_id"`
	PlanVersion string `bson:"plan_version"`
	UserId      string `bson:"user_id"`
	Position    struct {
		Unit     int `bson:"unit"`
		Exercise int `bson:"exercise"`
	} `bson:"position"`
	Data map[string]int `bson:"data"`
}
