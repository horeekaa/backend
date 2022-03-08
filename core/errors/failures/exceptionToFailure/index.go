package horeekaacoreexceptiontofailure

import (
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

	horeekaacoreexceptionenums.StoringImageFailed:              true,
	horeekaacoreexceptionenums.ClosingImageStoringWriterFailed: true,
	horeekaacoreexceptionenums.DeleteImageFailed:               true,

	horeekaacoreexceptionenums.ReverseGeocodeFailed: true,
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
func ConvertException(path string, err error) *horeekaacorebasefailure.Failure {
	if exception, ok := err.(*horeekaacorebaseexception.Exception); ok {
		if authenticationFailure[exception.Code] {
			return horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.AuthenticationTokenFailed,
				path,
				exception,
			)
		}

		if upstreamDuplicateFailure[exception.Code] {
			return horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.DuplicateObjectExist,
				path,
				exception,
			)
		}

		if objectNotFoundFailure[exception.Code] {
			return horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.ObjectNotFound,
				path,
				exception,
			)
		}

		if upstreamFailure[exception.Code] {
			return horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.UpstreamFailures,
				path,
				exception,
			)
		}
	}

	return horeekaacorefailure.NewFailureObject(
		horeekaacorefailureenums.UnknownFailures,
		path,
		err,
	)
}
