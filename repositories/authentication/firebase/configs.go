package firebaseauthentication

import (
	configs "github.com/horeekaa/backend/_commons/configs"
)

func GetFirebaseActionCodeSettings() map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"URL":                configs.GetEnvVariable(configs.FirebaseEmailActionCodeURL),
			"HandleCodeInApp":    configs.GetEnvVariable(configs.FirebaseEmailActionCodeHandleCodeInApp) == "true",
			"AndroidPackageName": configs.GetEnvVariable(configs.FirebaseEmailActionCodeAndroidPackageName),
			"AndroidInstallApp":  configs.GetEnvVariable(configs.FirebaseEmailActionCodeAndroidInstallApp) == "true",
		},
	}
}
