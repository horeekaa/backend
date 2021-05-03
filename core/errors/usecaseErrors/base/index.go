package horeekaacorebaseerror

import (
	"fmt"

	horeekaacorebasefailure "github.com/horeekaa/backend/core/errors/serviceFailures/base"
)

// Error struct for shaping usecase layer error
type Error struct {
	Message    string
	StatusCode int
	Path       string
	Err        *horeekaacorebasefailure.Failure
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error: %s, at %s. Details: %v", e.Message, e.Path, e.Err)
}
