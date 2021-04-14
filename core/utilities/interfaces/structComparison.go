package coreutilityinterfaces

type StructComparisonUtility interface {
	SetComparisonFunc(comparisonFunc func(tag interface{}, field1 interface{}, field2 interface{}, output *interface{})) (bool, error)
	SetPreDeepComparisonFunc(preDeepComparisonFunc func(tag interface{}, output *interface{})) (bool, error)
	CompareStructs(
		item1 interface{},
		item2 interface{},
		output *interface{},
	)
}
