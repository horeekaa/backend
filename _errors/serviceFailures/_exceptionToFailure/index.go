package horeekaaexceptiontofailure

import (
	horeekaabaseexception "github.com/horeekaa/backend/_errors/repoExceptions/_base"
	horeekaaexceptionenums "github.com/horeekaa/backend/_errors/repoExceptions/_enums"
	horeekaafailure "github.com/horeekaa/backend/_errors/serviceFailures"
	horeekaabasefailure "github.com/horeekaa/backend/_errors/serviceFailures/_base"
	horeekaafailureenums "github.com/horeekaa/backend/_errors/serviceFailures/_enums"
)

var upstreamFailure = map[string]bool{
	horeekaaexceptionenums.QueryObjectFailed:  true,
	horeekaaexceptionenums.CreateObjectFailed: true,
	horeekaaexceptionenums.UpdateObjectFailed: true,
	horeekaaexceptionenums.DeleteObjectFailed: true,
	horeekaaexceptionenums.SetAuthDataFailed:  true,
	horeekaaexceptionenums.DBConnectionFailed: true,
}

var upstreamDuplicateFailure = map[string]bool{
	horeekaaexceptionenums.DuplicateDataCreationNotAllowed: true,
}

var authenticationFailure = map[string]bool{
	horeekaaexceptionenums.BearerAuthTokenExpected: true,
	horeekaaexceptionenums.DecodingAuthTokenFailed: true,
}

var objectNotFoundFailure = map[string]bool{
	horeekaaexceptionenums.IDNotFound:        true,
	horeekaaexceptionenums.GetAuthDataFailed: true,
}

// ConvertException helps convert exceptions coming from the repo layer to be a failure in service layer
func ConvertException(path string, exception *horeekaabaseexception.Exception) *horeekaabasefailure.Failure {
	if authenticationFailure[exception.Message] {
		return horeekaafailure.NewFailureObject(
			horeekaafailureenums.AuthenticationTokenFailed,
			path,
			exception,
		)
	}

	if upstreamDuplicateFailure[exception.Message] {
		return horeekaafailure.NewFailureObject(
			horeekaafailureenums.DuplicateObjectExist,
			path,
			exception,
		)
	}

	if objectNotFoundFailure[exception.Message] {
		return horeekaafailure.NewFailureObject(
			horeekaafailureenums.ObjectNotFound,
			path,
			exception,
		)
	}

	if upstreamFailure[exception.Message] {
		return horeekaafailure.NewFailureObject(
			horeekaafailureenums.UpstreamFailures,
			path,
			exception,
		)
	}

	return horeekaafailure.NewFailureObject(
		horeekaafailureenums.UnknownFailures,
		path,
		exception,
	)
}
