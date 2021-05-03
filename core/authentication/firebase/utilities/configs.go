package firebaseauthcoreutilities

import (
	coreconfigs "github.com/horeekaa/backend/core/commons/configs"
)

func GetFirebaseActionCodeSettings() map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"URL":                coreconfigs.GetEnvVariable(coreconfigs.FirebaseEmailActionCodeURL),
			"HandleCodeInApp":    coreconfigs.GetEnvVariable(coreconfigs.FirebaseEmailActionCodeHandleCodeInApp) == "true",
			"AndroidPackageName": coreconfigs.GetEnvVariable(coreconfigs.FirebaseEmailActionCodeAndroidPackageName),
			"AndroidInstallApp":  coreconfigs.GetEnvVariable(coreconfigs.FirebaseEmailActionCodeAndroidInstallApp) == "true",
		},
	}
}
