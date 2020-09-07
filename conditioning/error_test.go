package conditioning

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck

	"fmt"
)

// ErrorEquals is a unit testing Checker verb for comparing errors.
var ErrorEquals Checker = &errorEqualsChecker{
	CheckerInfo{Name: "ErrorEquals", Params: []string{"obtained", "expected"}},
}

// errorEqualsChecker is a unit testing Checker verb for comparing errors.
type errorEqualsChecker struct {
	info CheckerInfo // A required bit of metadata for the checker.
}

// Info returns the checker metadata.
func (c *errorEqualsChecker) Info() (info *CheckerInfo) {
	return &c.info
}

// Checks does an actual unit test check.
func (c *errorEqualsChecker) Check(params []interface{}, names []string) (result bool, errstr string) {

	// If we obtained a no error, we always fail.
	if params[0] == nil {
		return false, ""
	}

	// The paramters must be our custom error and a string.
	obtained := params[0].(*errorWithStack)
	expected := params[1].(string)

	// If the message doesn't match there is an error.
	if obtained.message != expected {
		return false, fmt.Sprintf(`obtained error  = "%s"`, obtained.message)
	}

	// This check passed!
	return true, ""
}
