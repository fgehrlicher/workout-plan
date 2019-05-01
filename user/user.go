package user

type User struct {
	ID              string
	AuthorizedPlans []string
	ActivePlans []*PlanPointer
}
