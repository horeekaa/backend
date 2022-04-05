package horeekaacorefailuretoerror

import (
	"strings"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacorebaseerror "github.com/horeekaa/backend/core/errors/errors/base"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
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
	horeekaacorefailureenums.AcceptInvitationNotAllowed:          true,
	horeekaacorefailureenums.InvalidAccountPersonCredential:      true,

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
func ConvertFailure(path string, err error) *horeekaacorebaseerror.Error {
	if failure, ok := err.(*horeekaacorebasefailure.Failure); ok {
		failureCode := strings.Split(failure.Code, ".")[0]
		if authenticationError[failureCode] {
			return horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				path,
				failure,
			)
		}

		if forbiddenError[failureCode] {
			return horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.ForbiddenError,
				path,
				failure,
			)
		}

		if resourceNotFoundError[failureCode] {
			return horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.ResourceNotFoundError,
				path,
				failure,
			)
		}

		if conflictWithCurrentError[failureCode] {
			return horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.ConflictWithCurrentError,
				path,
				failure,
			)
		}

		if invalidInputError[failureCode] {
			return horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.InvalidInputError,
				path,
				failure,
			)
		}

		if badGatewayError[failureCode] {
			return horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.BadGatewayError,
				path,
				failure,
			)
		}
	}

	return horeekaacoreerror.NewErrorObject(
		horeekaacoreerrorenums.GeneralError,
		path,
		err,
	)
}
