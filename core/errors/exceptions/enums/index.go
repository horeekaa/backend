package horeekaacoreexceptionenums

// string enums for repo exception
const (
	BearerAuthTokenExpected string = "BEARER_AUTH_TOKEN_EXPECTED"
	DecodingAuthTokenFailed string = "DECODING_AUTH_TOKEN_FAILED"
	GetAuthDataFailed       string = "GET_AUTH_DATA_FAILED"
	SetAuthDataFailed       string = "SET_AUTH_DATA_FAILED"
	UpstreamException       string = "UPSTREAM_EXCEPTION"

	CreateObjectFailed              string = "CREATE_OBJECT_FAILED"
	DuplicateDataCreationNotAllowed string = "DUPLICATE_DATA_CREATION_NOT_ALLOWED"
	UpdateObjectFailed              string = "UPDATE_OBJECT_FAILED"
	DeleteObjectFailed              string = "DELETE_OBJECT_FAILED"
	QueryObjectFailed               string = "QUERY_OBJECT_FAILED"
	IDNotFound                      string = "ID_NOT_FOUND"
	DBConnectionFailed              string = "DB_CONNECTION_FAILED"

	StoringImageFailed              string = "STORING_IMAGE_FAILED"
	ClosingImageStoringWriterFailed string = "CLOSING_IMAGE_STORING_WRITER_FAILED"
	DeleteImageFailed               string = "DELETE_IMAGE_FAILED"

	SendNotifMessageFailed string = "SEND_NOTIF_MESSAGE_FAILED"

	ClientInitializationFailed string = "CLIENT_INITIALIZATION_FAILED"
)
