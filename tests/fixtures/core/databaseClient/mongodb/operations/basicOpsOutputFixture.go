package mongodbbasicoperationfixtures

import (
	"time"

	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var BasicOpsSingleResultOutput model.Account = model.Account{
	ID:        primitive.ObjectID{byte(255)},
	Email:     "accountone@test.com",
	Status:    "ACTIVE",
	Type:      "PERSON",
	CreatedAt: time.Now(),
	Person:    &model.Person{ID: primitive.ObjectID{byte(200)}},
}
