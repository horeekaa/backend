package horeekaabaseerror

import (
	"fmt"

	horeekaabasefailure "github.com/horeekaa/backend/_errors/serviceFailures/_base"
)

// Error struct for shaping usecase layer error
type Error struct {
	Message    string
	StatusCode int
	Path       string
	Err        *horeekaabasefailure.Failure
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error: %s, at %s. Details: %v", e.Message, e.Path, e.Err)
}
