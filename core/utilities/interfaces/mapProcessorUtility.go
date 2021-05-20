package coreutilityinterfaces

type MapProcessorUtility interface {
	RemoveNil(input map[string]interface{}) (bool, error)
	FlattenMap(
		keyPrefix string,
		input map[string]interface{},
		output *map[string]interface{},
	) (bool, error)
}
