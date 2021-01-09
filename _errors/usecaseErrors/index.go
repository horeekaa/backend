package horeekaaerror

import (
	horeekaabasefailures "github.com/horeekaa/backend/_errors/serviceFailures/_base"
	horeekaabaseerror "github.com/horeekaa/backend/_errors/usecaseErrors/_base"
)

// NewErrorObject getter usecaes layer Error Object
func NewErrorObject(message string, statusCode int, path string, err *horeekaabasefailures.Failure) *horeekaabaseerror.Error {
	return &horeekaabaseerror.Error{
		Message:    message,
		StatusCode: statusCode,
		Path:       path + (*err).Path,
		Err:        err,
	}
}
