package horeekaacoreerror

import (
	horeekaacorebaseerror "github.com/horeekaa/backend/core/errors/errors/base"
	horeekaacorebasefailure "github.com/horeekaa/backend/core/errors/failures/base"
)

// NewErrorObject getter usecaes layer Error Object
func NewErrorObject(message string, statusCode int, path string, err *horeekaacorebasefailure.Failure) *horeekaacorebaseerror.Error {
	extPath := ""
	if err != nil {
		if &err.Path != nil {
			extPath = err.Path
		}
	}

	return &horeekaacorebaseerror.Error{
		Message:    message,
		StatusCode: statusCode,
		Path:       path + extPath,
		Err:        err,
	}
}
