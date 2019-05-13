package template

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/knetic/govaluate"
)

type Evaluator struct {
	data                     map[string]interface{}
	enclosingStartCharacters string
	enclosingEndCharacters   string
}

func EvaluateTemplate(element interface{}, data map[string]int) error {
	rawData := make(map[string]interface{}, len(data))
	for key, value := range data {
		rawData[key] = value
	}

	evaluator := Evaluator{
		data:                     rawData,
		enclosingStartCharacters: "{{",
		enclosingEndCharacters:   "}}",
	}
	element = evaluator.resolve(element)
	return nil
}

func (evaluator *Evaluator) resolve(element interface{}) error {
	kind := reflect.Indirect(reflect.ValueOf(element)).Kind()
	value := reflect.ValueOf(element).Elem()
	return evaluator.resolveForKind(kind, &value)
}

func (evaluator *Evaluator) resolveForKind(kind reflect.Kind, value *reflect.Value) error {
	switch kind {
	case reflect.Struct:
		return evaluator.resolveStruct(value)
	case reflect.Slice:
		return evaluator.resolveSlice(value)
	case reflect.String:
		return evaluator.resolveString(value)
	}
	return nil
}

func (evaluator *Evaluator) resolveStruct(value *reflect.Value) error {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.IsValid() && field.CanSet() {
			kind := field.Kind()
			err := evaluator.resolveForKind(kind, &field)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (evaluator *Evaluator) resolveSlice(value *reflect.Value) error {
	for i := 0; i < value.Len(); i++ {
		sliceElement := value.Index(i)
		kind := sliceElement.Kind()
		err := evaluator.resolveForKind(kind, &sliceElement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (evaluator *Evaluator) resolveString(value *reflect.Value) error {
	if !value.CanSet() {
		return nil
	}

	valueString := value.String()
	if valueString == "" {
		return nil
	}

	startIndex := strings.Index(valueString, evaluator.enclosingStartCharacters)
	if startIndex == -1 {
		return nil
	}

	endIndex := strings.Index(valueString, evaluator.enclosingEndCharacters)
	if endIndex == -1 || endIndex < startIndex {
		return nil
	}

	evaluationString := valueString[startIndex+len(evaluator.enclosingStartCharacters) : endIndex]

	expression, err := govaluate.NewEvaluableExpression(evaluationString)
	if err != nil {
		return err
	}

	result, err := expression.Evaluate(evaluator.data)
	if err != nil {
		return err
	}

	newValueString := valueString[:startIndex]

	switch reflect.ValueOf(result).Kind() {
	case reflect.Float64:
		newValueString += fmt.Sprintf("%.2f", result)
	case reflect.Bool:
		newValueString += fmt.Sprintf("%v", result)
	}

	newValueString += valueString[endIndex+len(evaluator.enclosingEndCharacters):]

	value.SetString(newValueString)
	return nil
}
