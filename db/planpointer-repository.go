package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"workout-plan/plan"
)

const (
	collectionName = "planpointers"
	planIdKey      = "plan_id"
	planVersionKey = "plan_version"
	userIdKey      = "user_id"
	positionKey    = "position"
	dataKey        = "data"
	unitKey        = "unit"
	exerciseKey    = "exercise"
	startedKey     = "started"
	movedKey       = "moved"
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

func (planPointerRepository *PlanPointerRepository) Insert(pointer plan.Pointer) (*mongo.InsertOneResult, error) {
	planPointerBson := bsonx.Doc{
		{planIdKey, bsonx.String(pointer.PlanId)},
		{planVersionKey, bsonx.String(pointer.PlanVersion)},
		{userIdKey, bsonx.String(pointer.UserId)},
		{positionKey, bsonx.Document(
			bsonx.Doc{
				{unitKey, bsonx.Int32(int32(pointer.Position.Unit))},
				{exerciseKey, bsonx.Int32(int32(pointer.Position.Exercise))},
			}),
		},
		{startedKey, bsonx.Time(pointer.Started)},
	}

	var dataElems bsonx.Doc
	if len(pointer.Data) > 0 {
		for key, data := range pointer.Data {
			dataElems = dataElems.Append(key, bsonx.Int32(int32(data)))
		}
		planPointerBson = planPointerBson.Append(dataKey, bsonx.Document(dataElems))
	}

	return planPointerRepository.collection.InsertOne(planPointerRepository.requestContext, planPointerBson)
}

func (planPointerRepository *PlanPointerRepository) Update(pointer plan.Pointer) error {
	filter := bsonx.Doc{
		{planIdKey, bsonx.String(pointer.PlanId)},
		{planVersionKey, bsonx.String(pointer.PlanVersion)},
		{userIdKey, bsonx.String(pointer.UserId)},
	}

	updateDoc := bsonx.Doc{
		{positionKey, bsonx.Document(
			bsonx.Doc{
				{unitKey, bsonx.Int32(int32(pointer.Position.Unit))},
				{exerciseKey, bsonx.Int32(int32(pointer.Position.Exercise))},
			}),
		},
		{movedKey, bsonx.Time(time.Now())},
	}

	var dataElems bsonx.Doc
	if len(pointer.Data) > 0 {
		for key, data := range pointer.Data {
			dataElems = dataElems.Append(key, bsonx.Int32(int32(data)))
		}
		updateDoc = updateDoc.Append(dataKey, bsonx.Document(dataElems))
	}

	update := bsonx.Doc{
		{"$set", bsonx.Document(updateDoc)},
	}

	_, err := planPointerRepository.collection.UpdateOne(planPointerRepository.requestContext, filter, update)
	return err
}

func (planPointerRepository *PlanPointerRepository) Delete(pointer plan.Pointer) error {
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

func (planPointerRepository *PlanPointerRepository) GetAll(userId string) ([]plan.Pointer, error) {
	var userPlanPointers []plan.Pointer

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
		planPointer := &plan.Pointer{}
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

func (planPointerRepository *PlanPointerRepository) GetByPlan(userId string, planId string) (plan.Pointer, error) {
	singleResult := planPointerRepository.collection.FindOne(
		planPointerRepository.requestContext, bsonx.Doc{
			{userIdKey, bsonx.String(userId)},
			{planIdKey, bsonx.String(planId)},
		},
	)

	planPointer := plan.Pointer{Data: make(map[string]int)}
	err := singleResult.Decode(&planPointer)
	if err == mongo.ErrNoDocuments {
		err = NoPlanFoundError
	}

	return planPointer, err
}
