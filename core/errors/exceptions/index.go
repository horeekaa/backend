package horeekaacoreexception

import (
	horeekaacorebaseexception "github.com/horeekaa/backend/core/errors/exceptions/base"
)

// NewExceptionObject getter repo layer Exception Object
func NewExceptionObject(message string, path string, err error) *horeekaacorebaseexception.Exception {
	return &horeekaacorebaseexception.Exception{
		Message: message,
		Path:    path,
		Err:     err,
	}
}
