package models

import (
	"errors"
	"fmt"
	"reflect"
)

func TypeNotAllowedError(element interface{}) error {
	name := reflect.TypeOf(element).Name()
	elementValue := reflect.ValueOf(element)
	elementType := elementValue.MapIndex(reflect.ValueOf("Type")).Interface().(string)

	return errors.New(
		fmt.Sprintf(
			"type field `%v` is not allowed for %T.\nFull element: %+v",
			elementType,
			name,
			element,
		),
	)
}

func TypeNotEmptyValidator(element interface{}) error {
	name := reflect.TypeOf(element).Name()
	elementValue := reflect.ValueOf(element)
	elementType := elementValue.MapIndex(reflect.ValueOf("Type")).Interface().(string)

	if elementType == "" {
		return errors.New(
			fmt.Sprintf(
				"type field musnt be empty for `%v`.\nFull element: %+v",
				name,
				element,
			),
		)
	}

	return nil
}
