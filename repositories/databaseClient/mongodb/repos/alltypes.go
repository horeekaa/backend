package mongorepos

const (
	DefaultValuesUpdateType string = "UPDATE"
	DefaultValuesCreateType string = "CREATE"
)

type defaultValuesOptions struct {
	DefaultValuesType string
}
