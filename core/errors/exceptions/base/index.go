package horeekaacorebaseexception

import "fmt"

// Exception struct for shaping datasource layer error
type Exception struct {
	Code string   `json:"Code"`
	Path []string `json:"Path"`
	Err  error    `json:"Err"`
}

func (e *Exception) Error() string {
	if e.Err == nil {
		return ""
	}
	return fmt.Sprintf("Upstream Details: %s", e.Err.Error())
}
