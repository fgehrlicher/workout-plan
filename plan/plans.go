package plan

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Plans struct {
	underlyingSlice []Plan
}

func (plans *Plans) Add(plan Plan) {
	logEntry := log.WithFields(log.Fields{
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

var plansSingleton *Plans
var plansOnce sync.Once

func GetPlansInstance() *Plans {
	plansOnce.Do(func() {
		plansSingleton = &Plans{}
	})
	return plansSingleton
}

func InitializePlans(planDirectory string) error {
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
		log.WithFields(log.Fields{
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
		return nil, err
	}

	err = unmarshalMethod([]byte(data), plan)
	return plan, err
}