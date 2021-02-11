package mongomarshaler

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MarshalObjectID(objID primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("\"%s\"", objID.Hex())))
	})
}

func UnmarshalObjectID(v interface{}) (primitive.ObjectID, error) {
	parsed, ok := v.(string)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("ObjectID must be a string")
	}

	objectID, err := primitive.ObjectIDFromHex(parsed)
	return objectID, err
}
