package horeekaaexception

import (
	horeekaabaseexception "github.com/horeekaa/backend/_errors/repoExceptions/_base"
)

// NewExceptionObject getter repo layer Exception Object
func NewExceptionObject(message string, path string, err error) *horeekaabaseexception.Exception {
	return &horeekaabaseexception.Exception{
		Message: message,
		Path:    path,
		Err:     err,
	}
}
