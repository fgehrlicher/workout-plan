package plan_pointer

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "planpointers"

func NewPlanPointerRepository(database *mongo.Database) *PlanPointerRepository {
	return &PlanPointerRepository{
		collection: database.Collection(collectionName),
	}
}

type PlanPointerRepository struct {
	collection *mongo.Collection
}
