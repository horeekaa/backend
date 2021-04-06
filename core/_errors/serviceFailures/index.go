package horeekaacorefailure

import (
	horeekaacorebasefailure "github.com/horeekaa/backend/core/_errors/serviceFailures/_base"
)

// NewFailureObject getter service layer Failure
func NewFailureObject(message string, path string, err error) *horeekaacorebasefailure.Failure {
	return &horeekaacorebasefailure.Failure{
		Message: message,
		Path:    path + (*err).Path,
		Err:     err,
	}
}
