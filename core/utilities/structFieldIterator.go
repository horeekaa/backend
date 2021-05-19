package coreutilities

import (
	"reflect"

	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type structFieldIteratorUtility struct {
	iterateFunc        func(tag interface{}, field interface{}, output *interface{})
	preDeepIterateFunc func(tag interface{}, output *interface{})
}

func NewStructFieldIteratorUtility() (coreutilityinterfaces.StructFieldIteratorUtility, error) {
	return &structFieldIteratorUtility{}, nil
}

func (structFldIter *structFieldIteratorUtility) SetIteratingFunc(
	iterateFunc func(tag interface{}, field interface{}, output *interface{}),
) (bool, error) {
	structFldIter.iterateFunc = iterateFunc
	return true, nil
}

func (structFldIter *structFieldIteratorUtility) SetPreDeepIterateFunc(
	preDeepIterateFunc func(tag interface{}, output *interface{}),
) (bool, error) {
	structFldIter.preDeepIterateFunc = preDeepIterateFunc
	return true, nil
}

func (structFieldIterator *structFieldIteratorUtility) IterateStruct(
	item interface{},
	output *interface{},
) {
	itemType := reflect.TypeOf(item)
	itemReflectValue := reflect.ValueOf(item)
	itemReflectValue = reflect.Indirect(itemReflectValue)

	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	for i := 0; i < itemType.NumField(); i++ {
		itemTag := itemType.Field(i).Tag.Get("json")
		itemField := itemReflectValue.Field(i)
		if reflect.ValueOf(itemField.Interface()) == reflect.Zero(reflect.TypeOf(itemField.Interface())) {
			continue
		}

		if itemField.Kind() == reflect.Ptr {
			itemField = itemField.Elem()
		}

		if itemTag != "" && itemTag != "-" {
			if itemType.Field(i).Type.Kind() == reflect.Struct {
				if structFieldIterator.preDeepIterateFunc != nil {
					structFieldIterator.preDeepIterateFunc(itemTag, output)
				}
				structFieldIterator.IterateStruct(itemField.Interface(), output)
			} else {
				if structFieldIterator.iterateFunc != nil {
					structFieldIterator.iterateFunc(itemTag, itemField.Interface(), output)
				}
			}
		}
	}
}
