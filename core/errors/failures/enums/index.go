package horeekaacorefailureenums

// string enums for service failures
const (
	AuthenticationTokenFailed           string = "AUTHENTICATION_TOKEN_FAILED"
	AccountIsNotActive                  string = "ACCOUNT_IS_NOT_ACTIVE"
	AccountNotAllowedToUpdateEmail      string = "ACCOUNT_NOT_ALLOWED_TO_UPDATE_EMAIL"
	AccountIDNeededToRetrievePersonData string = "ACCOUNT_ID_NEEDED_TO_RETRIEVE_PERSON_DATA"
	FeatureNotAccessibleByAccount       string = "FEATURE_NOT_ACCESSIBLE_BY_ACCOUNT"
	SendEmailTypeNotExist               string = "SEND_EMAIL_TYPE_NOT_EXIST"
	AcceptInvitationNotAllowed          string = "ACCEPT_INVITATION_NOT_ALLOWED"

	MemberAccessRefNotExist                                   string = "MEMBER_ACCESS_REF_NOT_EXIST"
	OrganizationIDNeededToCreateOrganizationBasedMemberAccess string = "ORGANIZATION_ID_NEEDED_TO_CREATE_ORGANIZATION_BASED_MEMBER_ACCESS"

	DuplicateObjectExist string = "DUPLICATE_OBJECT_EXIST"
	ObjectNotFound       string = "OBJECT_NOT_FOUND"
	NothingToBeApproved  string = "NOTHING_TO_BE_APPROVED"

	POItemMismatchWithPOType            string = "PO_ITEM_MISMATCH_WITH_PO_TYPE"
	POSalesAmountExceedCreditLimit      string = "PO_SALES_AMOUNT_EXCEED_CREDIT_LIMIT"
	POMinimumOrderValueHasNotMet        string = "PO_MINIMUM_ORDER_VALUE_HAS_NOT_MET"
	POReturnAmountExceedFulfilledAmount string = "PO_RETURN_AMOUNT_EXCEED_FULFILLED_AMOUNT"
	UnapprovedPONotAllowedToFulfill     string = "UNAPPROVED_PO_NOT_ALLOWED_TO_FULFILL"
	SOReturnAmountExceedFulfilledAmount string = "SO_RETURN_AMOUNT_EXCEED_FULFILLED_AMOUNT"

	UpstreamFailures string = "UPSTREAM_FAILURES"
	UnknownFailures  string = "UNKNOWN_FAILURES"
)
