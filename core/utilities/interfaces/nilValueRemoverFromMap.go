package coreutilityinterfaces

type NilValueRemoverFromMapUtility interface {
	Execute(input map[string]interface{}) (bool, error)
}
