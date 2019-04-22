package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"workout-plan/models"

	"gopkg.in/yaml.v2"
)

func main() {
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
			switch filepath.Ext(fileInfo.Name()) {
			case ".yml":
				plan, err := loadYamlPlan(dirname + string(filepath.Separator) + fileInfo.Name())
				if err != nil {
					panic(err)
				}
				plans.Add(plan)
			case ".json":
				plan, err := loadJsonPlan(dirname + string(filepath.Separator) + fileInfo.Name())
				if err != nil {
					panic(err)
				}
				plans.Add(plan)
			}
		}
	}
	fmt.Printf("%+v\n", plans)
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
