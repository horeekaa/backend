package mongodbcorequerybuilderinterfaces

type MongoQueryBuilder interface {
	Execute(keyPrefix string, input interface{}, output *map[string]interface{}) (bool, error)
}
