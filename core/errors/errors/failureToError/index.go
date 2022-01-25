package horeekaacorefailuretoerror

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacorebaseerror "github.com/horeekaa/backend/core/errors/errors/base"
	horeekaacorebasefailure "github.com/horeekaa/backend/core/errors/failures/base"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
)

var authenticationError = map[string]bool{
	horeekaacorefailureenums.AuthenticationTokenFailed: true,
}

var forbiddenError = map[string]bool{
	horeekaacorefailureenums.AccountIsNotActive:             true,
	horeekaacorefailureenums.AccountNotAllowedToUpdateEmail: true,
	horeekaacorefailureenums.FeatureNotAccessibleByAccount:  true,
}

var resourceNotFoundError = map[string]bool{
	horeekaacorefailureenums.ObjectNotFound: true,
}

var conflictWithCurrentError = map[string]bool{
	horeekaacorefailureenums.DuplicateObjectExist: true,
}

var invalidInputError = map[string]bool{
	horeekaacorefailureenums.SendEmailTypeNotExist:               true,
	horeekaacorefailureenums.AccountIDNeededToRetrievePersonData: true,

	horeekaacorefailureenums.MemberAccessRefNotExist:                                   true,
	horeekaacorefailureenums.OrganizationIDNeededToCreateOrganizationBasedMemberAccess: true,

	horeekaacorefailureenums.NothingToBeApproved: true,

	horeekaacorefailureenums.POItemMismatchWithPOType:            true,
	horeekaacorefailureenums.POSalesAmountExceedCreditLimit:      true,
	horeekaacorefailureenums.POMinimumOrderValueHasNotMet:        true,
	horeekaacorefailureenums.POReturnAmountExceedFulfilledAmount: true,
	horeekaacorefailureenums.UnapprovedPONotAllowedToFulfill:     true,

	horeekaacorefailureenums.SOReturnAmountExceedFulfilledAmount: true,
}

var generalError = map[string]bool{
	horeekaacorefailureenums.UnknownFailures: true,
}

var badGatewayError = map[string]bool{
	horeekaacorefailureenums.UpstreamFailures: true,
}

// ConvertFailure helps convert failures coming from the service layer to be an error in usecase layer
func ConvertFailure(path string, errorObject interface{}) *horeekaacorebaseerror.Error {
	var failure horeekaacorebasefailure.Failure
	jsonTemp, _ := json.Marshal(errorObject)
	json.Unmarshal(jsonTemp, &failure)

	errMsg := ""
	if &failure.Message != nil {
		errMsg = failure.Message
	}

	if authenticationError[errMsg] {
		return horeekaacoreerror.NewErrorObject(
			failure.Message,
			401,
			path,
			&failure,
		)
	}

	if forbiddenError[errMsg] {
		return horeekaacoreerror.NewErrorObject(
			failure.Message,
			403,
			path,
			&failure,
		)
	}

	if resourceNotFoundError[errMsg] {
		return horeekaacoreerror.NewErrorObject(
			failure.Message,
			404,
			path,
			&failure,
		)
	}

	if conflictWithCurrentError[errMsg] {
		return horeekaacoreerror.NewErrorObject(
			failure.Message,
			409,
			path,
			&failure,
		)
	}

	if invalidInputError[errMsg] {
		return horeekaacoreerror.NewErrorObject(
			failure.Message,
			422,
			path,
			&failure,
		)
	}

	if badGatewayError[errMsg] {
		return horeekaacoreerror.NewErrorObject(
			failure.Message,
			502,
			path,
			&failure,
		)
	}

	return horeekaacoreerror.NewErrorObject(
		failure.Message,
		500,
		path,
		&failure,
	)
}
