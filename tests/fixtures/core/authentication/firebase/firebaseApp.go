package firebaseauthcoreclientfixtures

import (
	"context"
	"errors"
)

type FirebaseAppAuthOKFixture struct{}

func (_ *FirebaseAppAuthOKFixture) Auth(ctx context.Context) (interface{}, error) {
	return new(interface{}), nil
}

type FirebaseAppAuthErrorFixture struct{}

func (_ *FirebaseAppAuthErrorFixture) Auth(ctx context.Context) (interface{}, error) {
	return nil, errors.New("Some Upstream Error")
}
