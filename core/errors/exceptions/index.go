package horeekaacoreexception

import (
	horeekaacorebaseexception "github.com/horeekaa/backend/core/errors/exceptions/base"
)

// NewExceptionObject getter repo layer Exception Object
func NewExceptionObject(code string, path string, err error) *horeekaacorebaseexception.Exception {
	return &horeekaacorebaseexception.Exception{
		Code: code,
		Path: []string{path},
		Err:  err,
	}
}
