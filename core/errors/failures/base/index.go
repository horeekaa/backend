package horeekaacorebasefailure

// Failure struct for shaping repository layer error
type Failure struct {
	Code string   `json:"Code"`
	Path []string `json:"Path"`
	Err  error    `json:"Err"`
}

func (f *Failure) Error() string {
	if f.Err == nil {
		return f.Code
	}
	return f.Err.Error()
}
