package horeekaacoreerror

import (
	horeekaacorebaseerror "github.com/horeekaa/backend/core/errors/errors/base"
)

// NewErrorObject getter usecaes layer Error Object
func NewErrorObject(message string, statusCode int, path string, err error) *horeekaacorebaseerror.Error {
	return &horeekaacorebaseerror.Error{
		Message:    message,
		StatusCode: statusCode,
		Path:       path + (*err).Path,
		Err:        err,
	}
}
