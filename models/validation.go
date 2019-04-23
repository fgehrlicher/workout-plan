package models

import (
	"errors"
	"fmt"
	"reflect"
)

func TypeNotAllowedError(element interface{}) error {
	name := reflect.TypeOf(element).Name()
	elementValue := reflect.ValueOf(element)
	elementType := elementValue.MapIndex(reflect.ValueOf("Type"))

	return errors.New(
		fmt.Sprintf(
			"type field `%v` is not allowed for %T.\nFull element: %+v",
			elementType,
			name,
			element,
		),
	)
}
