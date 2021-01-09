package horeekaabaseexception

import "fmt"

// Exception struct for shaping repo layer error
type Exception struct {
	Message string
	Path    string
	Err     error
}

func (e *Exception) Error() string {
	return fmt.Sprintf("Error: %s, at %s. Details: %v", e.Message, e.Path, e.Err)
}
