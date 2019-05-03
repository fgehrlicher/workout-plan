package plan_pointer

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const collectionName = "planpointers"

func NewPlanPointerRepository(database *mongo.Database) (*PlanPointerRepository, error) {
	repository := &PlanPointerRepository{
		collection: database.Collection(collectionName),
	}

	return repository, repository.init()
}

type PlanPointerRepository struct {
	collection *mongo.Collection
}

func (planPointerRepository *PlanPointerRepository) init() error {
	indexView := planPointerRepository.collection.Indexes()

	planPointersIndex := mongo.IndexModel{
		Keys: bsonx.Doc{
			{"plan_id", bsonx.Int32(1)},
			{"plan_version", bsonx.Int32(1)},
			{"user_id", bsonx.Int32(1)},
		},
		Options: options.Index().SetName("plan-version-user"),
	}

	name, err := indexView.CreateOne(
		context.Background(),
		planPointersIndex,
	)

	fmt.Println(name)
	return err
}

func (planPointerRepository *PlanPointerRepository) Insert(pointer *PlanPointer) (*mongo.InsertOneResult, error) {
	planPointerBson := bson.D{
		{"plan_id", pointer.PlanId},
		{"plan_version", pointer.PlanVersion},
		{"user_id", pointer.UserId},
		{"position", bson.D{
			{"unit_id", pointer.Position.Unit.Id},
			{"exercise_key", pointer.Position.ExerciseKey},
		}},
	}

	return planPointerRepository.collection.InsertOne(context.Background(), planPointerBson)
}
