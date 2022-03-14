package horeekaacorebaseerror

import (
	"fmt"
	"strings"
)

// Error struct for shaping usecase layer error
type Error struct {
	Code string   `json:"Code"`
	Path []string `json:"Path"`
	Err  error    `json:"Err"`
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Code
	}
	return fmt.Sprintf("%s, at %s, details: %s", e.Code, strings.Join(e.Path, "/"), e.Err.Error())
}
