package horeekaacoreexceptiontofailure

import (
	"encoding/json"

	horeekaacorebaseexception "github.com/horeekaa/backend/core/errors/exceptions/base"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorebasefailure "github.com/horeekaa/backend/core/errors/failures/base"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
)

var upstreamFailure = map[string]bool{
	horeekaacoreexceptionenums.QueryObjectFailed:  true,
	horeekaacoreexceptionenums.CreateObjectFailed: true,
	horeekaacoreexceptionenums.UpdateObjectFailed: true,
	horeekaacoreexceptionenums.DeleteObjectFailed: true,
	horeekaacoreexceptionenums.SetAuthDataFailed:  true,
	horeekaacoreexceptionenums.DBConnectionFailed: true,
	horeekaacoreexceptionenums.UpstreamException:  true,

	horeekaacoreexceptionenums.ClientInitializationFailed: true,
}

var upstreamDuplicateFailure = map[string]bool{
	horeekaacoreexceptionenums.DuplicateDataCreationNotAllowed: true,
}

var authenticationFailure = map[string]bool{
	horeekaacoreexceptionenums.BearerAuthTokenExpected: true,
	horeekaacoreexceptionenums.DecodingAuthTokenFailed: true,
}

var objectNotFoundFailure = map[string]bool{
	horeekaacoreexceptionenums.IDNotFound:        true,
	horeekaacoreexceptionenums.GetAuthDataFailed: true,
}

// ConvertException helps convert exceptions coming from the repo layer to be a failure in service layer
func ConvertException(path string, errorObject interface{}) *horeekaacorebasefailure.Failure {
	var exception *horeekaacorebaseexception.Exception
	jsonTemp, _ := json.Marshal(errorObject)
	json.Unmarshal(jsonTemp, exception)
	errMsg := ""
	if &exception.Message != nil {
		errMsg = exception.Message
	}

	if authenticationFailure[errMsg] {
		return horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AuthenticationTokenFailed,
			path,
			exception,
		)
	}

	if upstreamDuplicateFailure[errMsg] {
		return horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.DuplicateObjectExist,
			path,
			exception,
		)
	}

	if objectNotFoundFailure[errMsg] {
		return horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.ObjectNotFound,
			path,
			exception,
		)
	}

	if upstreamFailure[errMsg] {
		return horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.UpstreamFailures,
			path,
			exception,
		)
	}

	return horeekaacorefailure.NewFailureObject(
		horeekaacorefailureenums.UnknownFailures,
		path,
		exception,
	)
}
