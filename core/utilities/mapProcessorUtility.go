package coreutilities

import (
	"fmt"
	"reflect"
	"strings"

	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type mapProcessorUtility struct{}

func NewMapProcessorUtility() (coreutilityinterfaces.MapProcessorUtility, error) {
	return &mapProcessorUtility{}, nil
}

func (mapProcessorUtil *mapProcessorUtility) RemoveNil(input map[string]interface{}) (bool, error) {
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
			mapProcessorUtil.RemoveNil(t)
		}
	}
	return true, nil
}

func (mapProcessorUtil *mapProcessorUtility) FlattenMap(
	keyPrefix string,
	input map[string]interface{},
	output *map[string]interface{},
) (bool, error) {
	val := reflect.ValueOf(input)
	prefix := keyPrefix
	if prefix != "" {
		prefix = fmt.Sprintf("%s.", prefix)
	}

	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)

		switch t := v.Interface().(type) {
		// If key is a JSON object (Go Map), use recursion to go deeper
		case map[string]interface{}:
			mapProcessorUtil.FlattenMap(
				fmt.Sprintf("%s%s", prefix, e.String()),
				t,
				output,
			)
		default:
			if strings.Contains(e.String(), "$") {
				if (*output)[keyPrefix] == nil {
					(*output)[keyPrefix] = map[string]interface{}{}
				}
				(*output)[keyPrefix].(map[string]interface{})[e.String()] = v.Interface()
				continue
			}
			(*output)[fmt.Sprintf("%s%s", prefix, e.String())] = v.Interface()

		}
	}
	return true, nil
}
