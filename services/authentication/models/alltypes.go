package authservicemodels

import (
	firebaseauthmodels "github.com/horeekaa/backend/repositories/authentication/firebase/models"
)

type UpdateAuthUserData struct {
	ServiceUpdateUser *firebaseauthmodels.UpdateAuthUserData
}
