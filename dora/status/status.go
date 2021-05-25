package status

import (
	"fmt"
)


type Status struct {
	Code Code
	Message string
}
func (s *Status) Error() string {
	return fmt.Sprintf("rpc error: code=%s, message:%s", s.Code, s.Message)
}

func New(code Code, message string) *Status {
	return &Status{
		Code:    code,
		Message: message,
	}
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(code Code, format string, a ...interface{}) *Status {
	return New(code, fmt.Sprintf(format, a...))
}

type Code uint32
const (
	Ok Code = 0
	Unimplemented Code = 1
	BadRequest Code = 2
	NotFound Code = 3
	PermissionDenied Code = 4
	Internal Code = 5
	Unauthenticated Code = 6
	Unknown Code = 7
	Unavailable Code = 8
)

var code2Str = map[Code]string{
	Ok:"Ok",
	Unimplemented: "Unimplemented",
	BadRequest: "BadRequest",
	NotFound: "NotFound",
	PermissionDenied: "PermissionDenied",
	Internal: "Internal",
	Unauthenticated: "Unauthenticated",
	Unknown: "Unknown",
	Unavailable: "Unavailable",
}
func (c Code) String() string {
	str, ok := code2Str[c]
	if !ok {
		return "Unknown"
	}

	return str
}