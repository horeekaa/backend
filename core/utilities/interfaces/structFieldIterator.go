package coreutilityinterfaces

type StructFieldIteratorUtility interface {
	SetIteratingFunc(iterateFunc func(tag interface{}, field interface{}, output *interface{})) (bool, error)
	SetPreDeepIterateFunc(preDeepIterate func(tag interface{}, output *interface{})) (bool, error)
	IterateStruct(
		item interface{},
		output *interface{},
	)
}
