package horeekaacorebasefailure

import (
	"fmt"

	horeekaacorebaseexception "github.com/horeekaa/backend/core/errors/exceptions/base"
)

// Failure struct for shaping service layer error
type Failure struct {
	Message string                               `json:"Message"`
	Path    string                               `json:"Path"`
	Err     *horeekaacorebaseexception.Exception `json:"Err"`
}

func (f *Failure) Error() string {
	return fmt.Sprintf("Error: %s, at %s. Details: %v", f.Message, f.Path, f.Err)
}
