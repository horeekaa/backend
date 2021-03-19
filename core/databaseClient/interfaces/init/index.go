package databasecoreclientinterfaces

type DatabaseClient interface {
	Connect() (bool, error)
}
