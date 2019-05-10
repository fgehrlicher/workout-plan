package plan_pointer

import (
	"context"
	"errors"
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
	unitKeyKey     = "unit_key"
	exerciseKeyKey = "exercise_key"
)

var (
	NoPlanFoundError = errors.New("no active plan for that id and user found")
)

func NewPlanPointerRepository(database *mongo.Database, requestTimeout time.Duration) *PlanPointerRepository {
	requestContext, _ := context.WithTimeout(context.Background(), requestTimeout)
	return &PlanPointerRepository{
		collection:     database.Collection(collectionName),
		requestContext: requestContext,
	}
}

type PlanPointerRepository struct {
	collection     *mongo.Collection
	requestContext context.Context
}

func (planPointerRepository *PlanPointerRepository) InitIndices() error {
	indexView := planPointerRepository.collection.Indexes()
	planPointersIndex := mongo.IndexModel{
		Keys: bsonx.Doc{
			{planIdKey, bsonx.Int32(1)},
			{userIdKey, bsonx.Int32(1)},
		},
		Options: options.Index().SetName("plan-version-user"),
	}

	_, err := indexView.CreateOne(
		planPointerRepository.requestContext,
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
				{unitKeyKey, bsonx.Int32(int32(pointer.Position.UnitKey))},
				{exerciseKeyKey, bsonx.Int32(int32(pointer.Position.ExerciseKey))},
			}),
		},
	}

	return planPointerRepository.collection.InsertOne(planPointerRepository.requestContext, planPointerBson)
}

func (planPointerRepository *PlanPointerRepository) Update(pointer PlanPointer) error {
	filter := bsonx.Doc{
		{planIdKey, bsonx.String(pointer.PlanId)},
		{planVersionKey, bsonx.String(pointer.PlanVersion)},
		{userIdKey, bsonx.String(pointer.UserId)},
	}

	update := bsonx.Doc{
		{positionKey, bsonx.Document(
			bsonx.Doc{
				{unitKeyKey, bsonx.Int32(int32(pointer.Position.UnitKey))},
				{exerciseKeyKey, bsonx.Int32(int32(pointer.Position.ExerciseKey))},
			}),
		},
	}

	_, err := planPointerRepository.collection.UpdateOne(planPointerRepository.requestContext, filter, update)
	return err
}

func (planPointerRepository *PlanPointerRepository) Delete(pointer PlanPointer) error {
	planPointerBson := bsonx.Doc{
		{planIdKey, bsonx.String(pointer.PlanId)},
		{planVersionKey, bsonx.String(pointer.PlanVersion)},
		{userIdKey, bsonx.String(pointer.UserId)},
	}

	singleResult := planPointerRepository.collection.FindOneAndDelete(
		planPointerRepository.requestContext,
		planPointerBson,
	)

	return singleResult.Decode(nil)
}

func (planPointerRepository *PlanPointerRepository) GetAll(userId string) ([]PlanPointer, error) {
	var userPlanPointers []PlanPointer

	cursor, err := planPointerRepository.collection.Find(
		planPointerRepository.requestContext, bsonx.Doc{
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

func (planPointerRepository *PlanPointerRepository) GetByPlan(userId string, planId string) (PlanPointer, error) {
	singleResult := planPointerRepository.collection.FindOne(
		planPointerRepository.requestContext, bsonx.Doc{
			{userIdKey, bsonx.String(userId)},
			{planIdKey, bsonx.String(planId)},
		},
	)

	planPointer := PlanPointer{}
	err := singleResult.Decode(&planPointer)
	if err == mongo.ErrNoDocuments {
		err = NoPlanFoundError
	}

	return planPointer, err
}
