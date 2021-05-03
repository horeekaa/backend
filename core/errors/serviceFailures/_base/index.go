package horeekaacorebasefailure

import (
	"fmt"

	horeekaacorebaseexception "github.com/horeekaa/backend/core/errors/repoExceptions/_base"
)

// Failure struct for shaping service layer error
type Failure struct {
	Message string
	Path    string
	Err     *horeekaacorebaseexception.Exception
}

func (f *Failure) Error() string {
	return fmt.Sprintf("Error: %s, at %s. Details: %v", f.Message, f.Path, f.Err)
}
