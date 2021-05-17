package coreutilities

import (
	"reflect"

	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type nilValueRemoverFromMapUtility struct{}

func NewNilValueRemoverFromMapUtility() (coreutilityinterfaces.NilValueRemoverFromMapUtility, error) {
	return &nilValueRemoverFromMapUtility{}, nil
}

func (nilValueRemover *nilValueRemoverFromMapUtility) Execute(input map[string]interface{}) (bool, error) {
	val := reflect.ValueOf(input)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(input, e.String())
			continue
		}
		switch t := v.Interface().(type) {
		// If key is a JSON object (Go Map), use recursion to go deeper
		case map[string]interface{}:
			nilValueRemover.Execute(t)
		}
	}
	return true, nil
}
