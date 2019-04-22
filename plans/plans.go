package plans

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"workout-plan/models"

	"gopkg.in/yaml.v2"
)

func InitializePlans() {
	dirname := "plans" + string(filepath.Separator)
	plans := models.Plans{}

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