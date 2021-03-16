package horeekaacoreerror

import (
	horeekaacorebasefailures "github.com/horeekaa/backend/core/_errors/serviceFailures/_base"
	horeekaacorebaseerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_base"
)

// NewErrorObject getter usecaes layer Error Object
func NewErrorObject(message string, statusCode int, path string, err *horeekaacorebasefailures.Failure) *horeekaacorebaseerror.Error {
	return &horeekaacorebaseerror.Error{
		Message:    message,
		StatusCode: statusCode,
		Path:       path + (*err).Path,
		Err:        err,
	}
}
