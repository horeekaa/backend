package mongodbcorequerybuilders

import (
	"fmt"
	"reflect"

	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	"github.com/horeekaa/backend/model"
)

type mongoQueryBuilder struct{}

func NewMongoQueryBuilder() (mongodbcorequerybuilderinterfaces.MongoQueryBuilder, error) {
	return &mongoQueryBuilder{}, nil
}

func (queryBuild *mongoQueryBuilder) Execute(
	keyPrefix string,
	input interface{},
	output *map[string]interface{},
) (bool, error) {
	refType := reflect.TypeOf(input)
	refValue := reflect.ValueOf(input)
	refValue = reflect.Indirect(refValue)
	prefix := keyPrefix
	if prefix != "" {
		prefix = fmt.Sprintf("%s.", prefix)
	}

	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
	}

	for i := 0; i < refType.NumField(); i++ {
		fieldName := refType.Field(i).Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			continue
		}

		refValueField := refValue.Field(i)
		if reflect.ValueOf(refValueField.Interface()) == reflect.Zero(refValueField.Type()) {
			continue
		}

		if refValueField.Kind() == reflect.Ptr {
			refValueField = refValueField.Elem()
		}

		switch refValueField.Kind() {
		case reflect.Struct:
			strValue, ok := refValueField.Interface().(model.StringFilterField)
			if ok {
				var result interface{}
				queryBuild.clientRequestToMongoQueryTranslation(
					strValue.Operation.String(),
					strValue.Value,
					strValue.Values,
					&result,
				)
				(*output)[fmt.Sprintf("%s%s", prefix, fieldName)] = result

				continue
			}

			intValue, ok := refValueField.Interface().(model.IntFilterField)
			if ok {
				var result interface{}
				queryBuild.clientRequestToMongoQueryTranslation(
					intValue.Operation.String(),
					intValue.Value,
					nil,
					&result,
				)
				(*output)[fmt.Sprintf("%s%s", prefix, fieldName)] = result

				continue
			}

			boolValue, ok := refValueField.Interface().(model.BooleanFilterField)
			if ok {
				var result interface{}
				queryBuild.clientRequestToMongoQueryTranslation(
					boolValue.Operation.String(),
					boolValue.Value,
					nil,
					&result,
				)
				(*output)[fmt.Sprintf("%s%s", prefix, fieldName)] = result

				continue
			}

			timeValue, ok := refValueField.Interface().(model.TimeFilterField)
			if ok {
				var result interface{}
				queryBuild.clientRequestToMongoQueryTranslation(
					timeValue.Operation.String(),
					timeValue.Value,
					nil,
					&result,
				)
				(*output)[fmt.Sprintf("%s%s", prefix, fieldName)] = result

				continue
			}

			objectIDValue, ok := refValueField.Interface().(model.ObjectIDFilterField)
			if ok {
				var result interface{}
				queryBuild.clientRequestToMongoQueryTranslation(
					objectIDValue.Operation.String(),
					objectIDValue.Value,
					objectIDValue.Values,
					&result,
				)
				(*output)[fmt.Sprintf("%s%s", prefix, fieldName)] = result

				continue
			}

			queryBuild.Execute(
				fmt.Sprintf("%s%s", prefix, fieldName),
				refValueField.Interface(),
				output,
			)
			break

		default:
			(*output)[fmt.Sprintf("%s%s", prefix, fieldName)] = refValueField.Interface()
		}
	}
	return true, nil
}

func (queryBuild *mongoQueryBuilder) clientRequestToMongoQueryTranslation(
	operation string,
	value interface{},
	values interface{},
	output *interface{},
) (bool, error) {
	mongoQueryMap := map[string]interface{}{}
	switch operation {
	case model.StringOperationEqual.String():
		if value == nil {
			mongoQueryMap["$exists"] = false
		} else {
			mongoQueryMap["$eq"] = value
		}
		break

	case model.StringOperationContains.String():
		mongoQueryMap["$regex"] = value
		mongoQueryMap["$options"] = "i"
		break

	case model.StringOperationIn.String():
		mongoQueryMap["$in"] = values
		break

	case model.StringOperationNotIn.String():
		mongoQueryMap["$nin"] = values
		break

	case model.StringOperationNotEqual.String():
		if value == nil {
			mongoQueryMap["$exists"] = true
		} else {
			mongoQueryMap["$ne"] = value
		}
		break

	case model.NumericOperationLessThan.String():
		mongoQueryMap["$lt"] = value
		break

	case model.NumericOperationLessThanOrEqual.String():
		mongoQueryMap["$lte"] = value
		break

	case model.NumericOperationMoreThan.String():
		mongoQueryMap["$gt"] = value
		break

	case model.NumericOperationMoreThanOrEqual.String():
		mongoQueryMap["$gte"] = value
		break

	case model.ObjectIDOperationHasValues.String():
		mongoQueryMap["$all"] = values
		break
	}

	*output = mongoQueryMap
	return true, nil
}
