package errorhandler

import (
    "fmt"
)

type ArgError struct {
	Arg int    `json:"arg"`
	Msg string `json:"msg"`
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("error with argument %d: %s", e.Arg, e.Msg)
}
