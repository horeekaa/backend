package horeekaacorebaseexception

import "fmt"

// Exception struct for shaping repo layer error
type Exception struct {
	Message string `json:"Message"`
	Path    string `json:"Path"`
	Err     error  `json:"Err"`
}

func (e *Exception) Error() string {
	return fmt.Sprintf("Error: %s, at %s. Details: %v", e.Message, e.Path, e.Err)
}
