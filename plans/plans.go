package plans

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"workout-plan/models"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Plans struct {
	underlyingSlice []*models.Plan
}

func (plans *Plans) Add(plan *models.Plan) {
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

var instance *Plans
var once sync.Once

func InitializePlans() {
	dirname := "plans" + string(filepath.Separator)
	plans := GetInstance()

	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	fileInfos, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.Mode().IsRegular() {
			filePath := dirname + string(filepath.Separator) + fileInfo.Name()
			switch filepath.Ext(fileInfo.Name()) {
			case ".yml":
				plan, err := loadYamlPlan(filePath)
				if err != nil {
					panic(err)
				}
				plans.Add(plan)
			case ".json":
				plan, err := loadJsonPlan(filePath)
				if err != nil {
					panic(err)
				}
				plans.Add(plan)
			}
		}
	}
}

func GetInstance() *Plans {
	once.Do(func() {
		instance = &Plans{}
	})
	return instance
}

func loadYamlPlan(path string) (*models.Plan, error) {
	return loadPlan(path, yaml.Unmarshal)
}

func loadJsonPlan(path string) (*models.Plan, error) {
	return loadPlan(path, json.Unmarshal)
}

func loadPlan(path string, unmarshalMethod func([]byte, interface{}) error) (*models.Plan, error) {
	var plan *models.Plan

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = unmarshalMethod([]byte(data), plan)
	return plan, err
}