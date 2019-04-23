package models

import (
	"errors"
	"fmt"
	"reflect"
)

func GetTypeField(element interface{}) (string, error) {
	var elementValue = reflect.ValueOf(element)
	kind := elementValue.Kind()
	if kind == reflect.Ptr {
		elementValue = reflect.Indirect(elementValue)
		kind = elementValue.Kind()
	}
	if kind != reflect.Struct {
		return "", errors.New(
			fmt.Sprintf("element is no struct. Got %T, value: %v",
				element,
				element,
			),
		)
	}

	return elementValue.FieldByName("Type").String(), nil
}

func TypeNotAllowedError(element interface{}) error {
	typeField, err := GetTypeField(element)
	if err != nil {
		return err
	}

	return errors.New(
		fmt.Sprintf(
			"type field `%v` is not allowed for %T.\nFull element: %+v",
			typeField,
			reflect.ValueOf(element).Type().Name(),
			element,
		),
	)
}

func TypeNotEmptyValidator(element interface{}) error {
	typeField, err := GetTypeField(element)
	if err != nil {
		return err
	}

	if typeField == "" {
		return errors.New(
			fmt.Sprintf(
				"type field musnt be empty for `%v`.\nFull element: %+v",
				reflect.ValueOf(element).Type().Name(),
				element,
			),
		)
	}

	return nil
}
