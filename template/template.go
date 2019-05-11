package template

import (
	"reflect"
)

type Evaluator struct {
	data map[string]string
}

func EvaluateTemplate(element interface{}, data map[string]string) error {
	evaluator := Evaluator{data: data}
	element = evaluator.resolve(element)
	return nil
}

func (evaluator *Evaluator) resolve(element interface{}) interface{} {
	kind := reflect.Indirect(reflect.ValueOf(element)).Kind()
	value := reflect.ValueOf(element).Elem()
	evaluator.resolveForKind(kind, &value)
	return element
}

func (evaluator *Evaluator) resolveForKind(kind reflect.Kind, value *reflect.Value) {
	switch kind {
	case reflect.Struct:
		evaluator.resolveStruct(value)
	case reflect.Slice:
		evaluator.resolveSlice(value)
	case reflect.String:
		evaluator.resolveString(value)
	}
}

func (evaluator *Evaluator) resolveStruct(value *reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.IsValid() && field.CanSet() {
			kind := field.Kind()
			evaluator.resolveForKind(kind, &field)
		}
	}
}

func (evaluator *Evaluator) resolveSlice(value *reflect.Value) {
	for i := 0; i < value.Len(); i++ {
		sliceElement := value.Index(i)
		kind := sliceElement.Kind()
		evaluator.resolveForKind(kind, &sliceElement)
	}
}

func (evaluator *Evaluator) resolveString(value *reflect.Value) {
	if value.CanSet() {
		value.SetString("TEST STRING")
	}
}
