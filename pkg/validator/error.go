package validator

import "fmt"

var ErrValidator *ErrorValidator

type ErrorValidator struct {
	field string
	tag   string
	param string
	value interface{}
}

func (m *ErrorValidator) Error() string {
	return fmt.Sprintf("%s does not meet the '%s[%s]' requirement with value '%v'", m.field, m.tag, m.param, m.value)
}
