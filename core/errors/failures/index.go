package horeekaacorefailure

import (
	"fmt"

	horeekaacorebaseexception "github.com/horeekaa/backend/core/errors/exceptions/base"
	horeekaacorebasefailure "github.com/horeekaa/backend/core/errors/failures/base"
)

// NewFailureObject getter service layer Failure
func NewFailureObject(code string, path string, err error) *horeekaacorebasefailure.Failure {
	errPath := []string{path}
	extendedCode := code

	if exception, ok := err.(*horeekaacorebaseexception.Exception); ok {
		errPath = append(errPath, exception.Path...)
		extendedCode = fmt.Sprintf("%s.%s", code, exception.Code)
	}

	return &horeekaacorebasefailure.Failure{
		Code: extendedCode,
		Path: errPath,
		Err:  err,
	}
}
