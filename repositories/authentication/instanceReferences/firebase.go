package authenticationinstancereferences

import (
	firebaseauthinterfaces "github.com/horeekaa/backend/repositories/authentication/firebase/interfaces"
)

type AuthenticationRepo struct {
	Instance *firebaseauthinterfaces.FirebaseAuthentication
}
