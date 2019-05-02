package plan_pointer

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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

func (planPointerRepository *PlanPointerRepository) Insert(pointer *PlanPointer) error {
	planPointerBson := bson.D{
		{"plan_id", pointer.PlanId},
		{"position", bson.D{
			{"unit_id",pointer.Position.Unit.Id},
			{"exercise_key",pointer.Position.ExerciseKey},
		}},
	}

	_, err := planPointerRepository.collection.InsertOne(context.Background(), planPointerBson)
	return err
}
