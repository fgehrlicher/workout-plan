package plan_pointer

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	collectionName = "planpointers"
	planIdKey      = "plan_id"
	planVersionKey = "plan_version"
	userIdKey      = "user_id"
	positionKey    = "position"
	unitIdKey      = "unit_id"
	exerciseKeyKey = "exercise_key"
)

var requestContext context.Context

func NewPlanPointerRepository(database *mongo.Database, requestTimeout time.Duration) *PlanPointerRepository {
	requestContext, _ = context.WithTimeout(context.Background(), requestTimeout)
	return &PlanPointerRepository{
		collection: database.Collection(collectionName),
	}
}

type PlanPointerRepository struct {
	collection *mongo.Collection
}

func (planPointerRepository *PlanPointerRepository) InitIndices() error {
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
		requestContext,
		planPointersIndex,
	)

	return err
}

func (planPointerRepository *PlanPointerRepository) Insert(pointer PlanPointer) (*mongo.InsertOneResult, error) {
	planPointerBson := bsonx.Doc{
		{planIdKey, bsonx.String(pointer.PlanId)},
		{planVersionKey, bsonx.String(pointer.PlanVersion)},
		{userIdKey, bsonx.String(pointer.UserId)},
		{positionKey, bsonx.Document(
			bsonx.Doc{
				{unitIdKey, bsonx.String(pointer.Position.UnitId)},
				{exerciseKeyKey, bsonx.Int32(int32(pointer.Position.ExerciseKey))},
			}),
		},
	}

	return planPointerRepository.collection.InsertOne(requestContext, planPointerBson)
}

func (planPointerRepository *PlanPointerRepository) GetAll(userId string) ([]PlanPointer, error) {
	var userPlanPointers []PlanPointer

	cursor, err := planPointerRepository.collection.Find(
		requestContext, bsonx.Doc{
			{userIdKey, bsonx.String(userId)},
		},
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		planPointer := &PlanPointer{}
		err = cursor.Decode(planPointer)
		if err != nil {
			return nil, err
		}

		userPlanPointers = append(userPlanPointers, *planPointer)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return userPlanPointers, nil
}
