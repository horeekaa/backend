package mongodbbasicoperationfixtures

import (
	"time"

	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var layout = "2006-01-02T15:04:05.000Z"
var str = "2014-11-12T11:45:26.371Z"
var t, err = time.Parse(layout, str)

var BasicOpsSingleResultOutput model.Account = model.Account{
	ID:        primitive.ObjectID{byte(255)},
	Email:     "accountone@test.com",
	Status:    "ACTIVE",
	Type:      "PERSON",
	CreatedAt: t,
	Person:    &model.Person{ID: primitive.ObjectID{byte(200)}},
}
