package template

import (
	"reflect"
)

type Evaluator struct {
	data map[string]string
}

func EvaluateTemplate(element interface{}, data map[string]string) error {
	evaluator := Evaluator{data: data}
	evaluator.resolveForType(element)
	return nil
}

func (evaluator *Evaluator) resolveForType(element interface{}) {
	switch reflect.TypeOf(element).Kind() {
	case reflect.Struct:
		evaluator.resolveSlice(element)
	case reflect.Slice:
		evaluator.resolveSlice(element)
	case reflect.String:
		evaluator.resolveString(element)
	}
}

func (evaluator *Evaluator) resolveStruct(element interface{}) {

}

func (evaluator *Evaluator) resolveSlice(element interface{}) {

}

func (evaluator *Evaluator) resolveString(element interface{}) {

}
