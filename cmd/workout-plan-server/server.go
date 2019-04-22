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
				data, err := ioutil.ReadFile(dirname + string(filepath.Separator) + fileInfo.Name())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				plan := models.Plan{}
				err = yaml.Unmarshal([]byte(data), &plan)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				plans.Add(plan)
			case ".json":
				data, err := ioutil.ReadFile(dirname + string(filepath.Separator) + fileInfo.Name())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				plan := models.Plan{}
				err = json.Unmarshal([]byte(data), &plan)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				plans.Add(plan)
			}
		}
	}

	fmt.Printf("%+v\n", plans)

}
