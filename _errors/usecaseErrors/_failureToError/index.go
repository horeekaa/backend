package horeekaafailuretoerror

import (
	horeekaafailureenums "github.com/horeekaa/backend/_errors/serviceFailures/_enums"
	horeekaaerror "github.com/horeekaa/backend/_errors/usecaseErrors"
	horeekaabaseerror "github.com/horeekaa/backend/_errors/usecaseErrors/_base"
)

var authenticationError = map[string]bool{
	horeekaafailureenums.AuthenticationTokenFailed: true,
}

var forbiddenError = map[string]bool{
	horeekaafailureenums.AccountIsNotActive:             true,
	horeekaafailureenums.AccountNotAllowedToUpdateEmail: true,
	horeekaafailureenums.FeatureNotAccessibleByAccount:  true,
}

var resourceNotFoundError = map[string]bool{
	horeekaafailureenums.ObjectNotFound: true,
}

var conflictWithCurrentError = map[string]bool{
	horeekaafailureenums.DuplicateObjectExist: true,
}

var invalidInputError = map[string]bool{
	horeekaafailureenums.SendEmailTypeNotExist: true,
}

var generalError = map[string]bool{
	horeekaafailureenums.UnknownFailures: true,
}

var badGatewayError = map[string]bool{
	horeekaafailureenums.UpstreamFailures: true,
}

// ConvertFailure helps convert failures coming from the service layer to be an error in usecase layer
func ConvertFailure(path string, failure *error) *horeekaabaseerror.Error {
	errMsg := ""
	if &failure.Message != nil {
		errMsg := *failure.Message
	}

	if authenticationError[errMsg] {
		return horeekaaerror.NewErrorObject(
			(*failure).Message,
			401,
			path,
			failure,
		)
	}

	if forbiddenError[errMsg] {
		return horeekaaerror.NewErrorObject(
			(*failure).Message,
			403,
			path,
			failure,
		)
	}

	if resourceNotFoundError[errMsg] {
		return horeekaaerror.NewErrorObject(
			(*failure).Message,
			404,
			path,
			failure,
		)
	}

	if conflictWithCurrentError[errMsg] {
		return horeekaaerror.NewErrorObject(
			(*failure).Message,
			409,
			path,
			failure,
		)
	}

	if invalidInputError[errMsg] {
		return horeekaaerror.NewErrorObject(
			(*failure).Message,
			422,
			path,
			failure,
		)
	}

	if badGatewayError[errMsg] {
		return horeekaaerror.NewErrorObject(
			(*failure).Message,
			503,
			path,
			failure,
		)
	}

	return horeekaaerror.NewErrorObject(
		(*failure).Message,
		500,
		path,
		failure,
	)
}
