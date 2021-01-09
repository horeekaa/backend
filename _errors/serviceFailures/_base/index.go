package horeekaabasefailure

import (
	"fmt"

	horeekaabaseexception "github.com/horeekaa/backend/_errors/repoExceptions/_base"
)

// Failure struct for shaping service layer error
type Failure struct {
	Message string
	Path    string
	Err     *horeekaabaseexception.Exception
}

func (f *Failure) Error() string {
	return fmt.Sprintf("Error: %s, at %s. Details: %v", f.Message, f.Path, f.Err)
}
