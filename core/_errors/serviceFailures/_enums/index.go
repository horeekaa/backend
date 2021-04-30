package horeekaacorefailureenums

// string enums for service failures
const (
	AuthenticationTokenFailed           string = "AUTHENTICATION_TOKEN_FAILED"
	AccountIsNotActive                  string = "ACCOUNT_IS_NOT_ACTIVE"
	AccountNotAllowedToUpdateEmail      string = "ACCOUNT_NOT_ALLOWED_TO_UPDATE_EMAIL"
	AccountIDNeededToRetrievePersonData string = "ACCOUNT_ID_NEEDED_TO_RETRIEVE_PERSON_DATA"
	FeatureNotAccessibleByAccount       string = "FEATURE_NOT_ACCESSIBLE_BY_ACCOUNT"
	SendEmailTypeNotExist               string = "SEND_EMAIL_TYPE_NOT_EXIST"

	MemberAccessRefNotExist                                   string = "MEMBER_ACCESS_REF_NOT_EXIST"
	OrganizationIDNeededToCreateOrganizationBasedMemberAccess string = "ORGANIZATION_ID_NEEDED_TO_CREATE_ORGANIZATION_BASED_MEMBER_ACCESS"

	ExistingObjectAndItsIDMustNotBeNilForUpdateActivity string = "EXISTING_OBJECT_AND_ITS_ID_MUST_NOT_BE_NIL_FOR_UPDATE_ACTIVITY"

	DuplicateObjectExist string = "DUPLICATE_OBJECT_EXIST"
	ObjectNotFound       string = "OBJECT_NOT_FOUND"

	UpstreamFailures string = "UPSTREAM_FAILURES"
	UnknownFailures  string = "UNKNOWN_FAILURES"
)
