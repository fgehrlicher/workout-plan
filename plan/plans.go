package plan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"workout-plan/version"
)

var plansLogger *logrus.Logger

type Plans struct {
	underlyingSlice []Plan
}

func (plans *Plans) Add(plan Plan) {
	logEntry := plansLogger.WithFields(logrus.Fields{
		"Id":      plan.ID,
		"Version": plan.Version,
	})

	for _, existingPlan := range plans.underlyingSlice {
		if existingPlan.ID == plan.ID && existingPlan.Version == plan.Version {
			logEntry.Error("Plan id - version combination already exists")
			return
		}
	}

	plans.underlyingSlice = append(plans.underlyingSlice, plan)
	logEntry.Info("Plan added")
}

func (plans *Plans) GetLatest(planId string) (*Plan, error) {
	var returnPlan *Plan

	for key, plan := range plans.underlyingSlice {
		if plan.ID != planId {
			continue
		}

		if returnPlan == nil {
			returnPlan = &plans.underlyingSlice[key]
			continue
		}

		isBigger, err := version.IsGreater(returnPlan.Version, plan.Version)
		if err != nil {
			return nil, err
		}

		if isBigger {
			returnPlan = &plans.underlyingSlice[key]
		}
	}

	if returnPlan == nil {
		return nil, fmt.Errorf(
			"no plan with id `%v` found",
			planId,
		)
	}

	return returnPlan, nil
}

func (plans *Plans) GetAll() []*Plan {
	var returnPlans []*Plan
	for key := range plans.underlyingSlice {
		returnPlans = append(returnPlans, &plans.underlyingSlice[key])
	}

	return returnPlans
}

func (plans *Plans) GetAllLatest() ([]*Plan, error) {
	var (
		returnPlans []*Plan
		found       bool
	)

	for key, plan := range plans.underlyingSlice {

		found = false
		for addedPlansKey, addedPlan := range returnPlans {
			if plan.ID != addedPlan.ID {
				continue
			}
			found = true

			isBigger, err := version.IsGreater(addedPlan.Version, plan.Version)
			if err != nil {
				return nil, err
			}

			if isBigger {
				returnPlans[addedPlansKey] = &plans.underlyingSlice[key]
			}
			break
		}

		if !found {
			returnPlans = append(returnPlans, &plans.underlyingSlice[key])
		}
	}

	return returnPlans, nil
}

func (plans *Plans) Get(planId string, version string) (*Plan, error) {
	for _, plan := range plans.underlyingSlice {
		if plan.ID == planId && plan.Version == version {
			return &plan, nil
		}
	}

	return nil, fmt.Errorf(
		"no plan with id `%v` and version `%v` found",
		planId,
		version,
	)
}

var plansSingleton *Plans
var plansOnce sync.Once

func GetPlansInstance() *Plans {
	plansOnce.Do(func() {
		plansSingleton = &Plans{}
	})
	return plansSingleton
}

func InitializePlans(planDirectory string, logger *logrus.Logger) error {
	plansLogger = logger
	plans := GetPlansInstance()

	fileInfos, err := ioutil.ReadDir(planDirectory)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		if !fileInfo.Mode().IsRegular() {
			continue
		}

		filePath := planDirectory + string(filepath.Separator) + fileInfo.Name()
		switch filepath.Ext(fileInfo.Name()) {
		case ".yml":
			plan, err := loadYamlPlan(filePath)
			if err != nil {
				return err
			}

			err = validatePlan(plan)
			if err != nil {
				return err
			}

			plans.Add(*plan)
		case ".json":
			plan, err := loadJsonPlan(filePath)
			if err != nil {
				return err
			}

			err = validatePlan(plan)
			if err != nil {
				return err
			}

			plans.Add(*plan)
		}
	}

	return nil
}

func validatePlan(plan *Plan) error {
	err := plan.Validate()
	if err != nil {
		plansLogger.WithFields(logrus.Fields{
			"Id": plan.ID,
		}).Error("plan validation error")
		return err
	}
	return nil
}

func loadYamlPlan(path string) (*Plan, error) {
	return loadPlan(path, yaml.Unmarshal)
}

func loadJsonPlan(path string) (*Plan, error) {
	return loadPlan(path, json.Unmarshal)
}

func loadPlan(path string, unmarshalMethod func([]byte, interface{}) error) (*Plan, error) {
	plan := &Plan{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"can´t load plan file (tried: '%v'): %v",
			path,
			err,
		)
	}

	err = unmarshalMethod([]byte(data), plan)
	return plan, err
}
