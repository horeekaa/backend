package horeekaacorefailure

import (
	horeekaacorebaseexception "github.com/horeekaa/backend/core/errors/exceptions/base"
	horeekaacorebasefailure "github.com/horeekaa/backend/core/errors/failures/base"
)

// NewFailureObject getter service layer Failure
func NewFailureObject(message string, path string, err *horeekaacorebaseexception.Exception) *horeekaacorebasefailure.Failure {
	extPath := ""
	if err != nil {
		if &err.Path != nil {
			extPath = err.Path
		}
	}

	return &horeekaacorebasefailure.Failure{
		Message: message,
		Path:    path + extPath,
		Err:     err,
	}
}
