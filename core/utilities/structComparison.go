package coreutilities

import (
	"reflect"

	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type structComparisonUtility struct {
	comparisonFunc     func(tag interface{}, field1 interface{}, field2 interface{}, output *interface{})
	preDeepCompareFunc func(tag interface{}, output *interface{})
}

func NewStructComparisonUtility() (coreutilityinterfaces.StructComparisonUtility, error) {
	return &structComparisonUtility{}, nil
}

func (strctComparisonUtility *structComparisonUtility) SetComparisonFunc(
	comparisonFunc func(tag interface{}, field1 interface{}, field2 interface{}, output *interface{}),
) (bool, error) {
	strctComparisonUtility.comparisonFunc = comparisonFunc
	return true, nil
}

func (strctComparisonUtility *structComparisonUtility) SetPreDeepComparisonFunc(
	preDeepComparisonFunc func(tag interface{}, output *interface{}),
) (bool, error) {
	strctComparisonUtility.preDeepCompareFunc = preDeepComparisonFunc
	return true, nil
}

func (strctComparisonUtility *structComparisonUtility) CompareStructs(
	item1 interface{},
	item2 interface{},
	output *interface{},
) {
	if item1 == nil || item2 == nil {
		return
	}
	item1Type := reflect.TypeOf(item1)
	item1ReflectValue := reflect.ValueOf(item1)
	item1ReflectValue = reflect.Indirect(item1ReflectValue)

	item2ReflectValue := reflect.ValueOf(item2)
	item2ReflectValue = reflect.Indirect(item2ReflectValue)

	if item1Type.Kind() == reflect.Ptr {
		item1Type = item1Type.Elem()
	}

	for i := 0; i < item1Type.NumField(); i++ {
		item1Tag := item1Type.Field(i).Tag.Get("json")
		item1Field := item1ReflectValue.Field(i).Interface()
		item2Field := item2ReflectValue.FieldByName(item1Tag).Interface()

		if item1Tag != "" && item1Tag != "-" {
			if item1Type.Field(i).Type.Kind() == reflect.Struct {
				if strctComparisonUtility.preDeepCompareFunc != nil {
					strctComparisonUtility.preDeepCompareFunc(item1Tag, output)
				}
				strctComparisonUtility.CompareStructs(item1Field, item2Field, output)
			} else {
				if strctComparisonUtility.comparisonFunc != nil {
					strctComparisonUtility.comparisonFunc(item1Tag, item1Field, item2Field, output)
				}
			}
		}
	}
}
