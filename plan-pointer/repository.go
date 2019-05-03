package plan_pointer

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	collectionName = "planpointers"
	planIdKey      = "plan_id"
	planVersionKey = "plan_version"
	userIdKey      = "user_id"
)

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
			{planIdKey, bsonx.Int32(1)},
			{planVersionKey, bsonx.Int32(1)},
			{userIdKey, bsonx.Int32(1)},
		},
		Options: options.Index().SetName("plan-version-user"),
	}

	_, err := indexView.CreateOne(
		context.Background(),
		planPointersIndex,
	)

	return err
}

func (planPointerRepository *PlanPointerRepository) Insert(pointer PlanPointer) (*mongo.InsertOneResult, error) {
	planPointerBson := bsonx.Doc{
		{planIdKey, bsonx.String(pointer.PlanId)},
		{planVersionKey, bsonx.String(pointer.PlanVersion)},
		{userIdKey, bsonx.String(pointer.UserId)},
		{"position", bsonx.Document(bsonx.Doc{
			{"unit_id", bsonx.String(pointer.Position.Unit.Id)},
			{"exercise_key", bsonx.Int32(int32(pointer.Position.ExerciseKey))},
		}),
		},
	}

	return planPointerRepository.collection.InsertOne(context.Background(), planPointerBson)
}
