package dependency

import "fmt"

type ValidationError struct {
	Msg    string
	Errors []error
}

func (v *ValidationError) Error() string {
	msg := v.Msg
	for _, err := range v.Errors {
		msg += fmt.Sprintf("\n  %s", err.Error())
	}

	return msg
}

type errorList []error

func (el *errorList) AddError(err error) {
	*el = append(*el, err)
}
