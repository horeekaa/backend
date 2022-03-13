package horeekaacoreerror

import (
	"fmt"

	horeekaacorebaseerror "github.com/horeekaa/backend/core/errors/errors/base"
	horeekaacorebasefailure "github.com/horeekaa/backend/core/errors/failures/base"
)

// NewErrorObject getter usecaes layer Error Object
func NewErrorObject(code string, path string, err error) *horeekaacorebaseerror.Error {
	errPath := []string{path}
	extendedCode := fmt.Sprintf("err.horeekaa.%s", code)

	if failure, ok := err.(*horeekaacorebasefailure.Failure); ok {
		errPath = append(errPath, failure.Path...)
		extendedCode = fmt.Sprintf("err.horeekaa.%s.%s", code, failure.Code)
	}

	return &horeekaacorebaseerror.Error{
		Code: extendedCode,
		Path: errPath,
		Err:  err,
	}
}
