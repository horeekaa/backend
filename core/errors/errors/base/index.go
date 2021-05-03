package horeekaacorebaseerror

import (
	"fmt"

	horeekaacorebasefailure "github.com/horeekaa/backend/core/errors/failures/base"
)

// Error struct for shaping usecase layer error
type Error struct {
	Message    string                           `json:"Message"`
	StatusCode int                              `json:"StatusCode"`
	Path       string                           `json:"Path"`
	Err        *horeekaacorebasefailure.Failure `json:"Err"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error: %s, at %s. Details: %v", e.Message, e.Path, e.Err)
}
