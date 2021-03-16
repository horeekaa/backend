package coreconfigs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// environment variable configuration keys
const (
	DbConfigURL     string = "DBCONFIG_URL"
	DbConfigDBName  string = "DBCONFIG_DBNAME"
	DbConfigTimeout string = "DBCONFIG_TIMEOUT"

	FirebaseConfigProjectID                   string = "FIREBASECONFIG_PROJECT_ID"
	FirebaseServiceAccountPath                string = "FIREBASECONFIG_SERVICE_ACCOUNT_PATH"
	FirebaseEmailActionCodeURL                string = "FIREBASECONFIG_EMAIL_ACTION_CODE_URL"
	FirebaseEmailActionCodeAndroidPackageName string = "FIREBASECONFIG_EMAIL_ACTION_CODE_ANDROID_PACKAGENAME"
	FirebaseEmailActionCodeAndroidInstallApp  string = "FIREBASECONFIG_EMAIL_ACTION_CODE_ANDROID_INSTALLAPP"
	FirebaseEmailActionCodeHandleCodeInApp    string = "FIREBASECONFIG_EMAIL_ACTION_CODE_HANDLECODEINAPP"

	GoogleCloudConfigStorageBucketName string = "GOOGLECLOUDCONFIG_STORAGE_BUCKET_NAME"

	SendGridConfigKey                   string = "SENDGRIDCONFIG_KEY"
	SendGridConfigTemplateResetPassword string = "SENDGRIDCONFIG_TEMPLATES_RESET_PASSWORD"
	SendGridConfigTemplateVerifyEmail   string = "SENDGRIDCONFIG_TEMPLATES_VERIFY_EMAIL"
	SendGridConfigDefaultEmailSender    string = "SENDGRIDCONFIG_DEFAULT_EMAIL_SENDER"

	SystemWideTimeFormat string = "SYSTEM_WIDE_TIME_FORMAT"
)

// GetEnvVariable will retrieve the value of the environment variable based on the input key
func GetEnvVariable(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Cannot load .env file!")
	}

	return os.Getenv(key)
}
