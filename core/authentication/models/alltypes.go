package authenticationcoremodels

type UpdateAuthUserData struct {
	UID           string
	Email         string
	EmailVerified bool
	PhoneNumber   string
	Password      string
	DisplayName   string
	PhotoURL      string
	Disabled      bool
}

type contextKey struct {
	name string
}

var UserContextKey *contextKey = &contextKey{"user"}
