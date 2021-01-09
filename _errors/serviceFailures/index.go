package horeekaafailure

import (
	horeekaabaseexception "github.com/horeekaa/backend/_errors/repoExceptions/_base"
	horeekaabasefailure "github.com/horeekaa/backend/_errors/serviceFailures/_base"
)

// NewFailureObject getter service layer Failure
func NewFailureObject(message string, path string, err *horeekaabaseexception.Exception) *horeekaabasefailure.Failure {
	return &horeekaabasefailure.Failure{
		Message: message,
		Path:    path + (*err).Path,
		Err:     err,
	}
}
