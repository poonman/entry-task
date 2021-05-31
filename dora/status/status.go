package status

import (
	"fmt"
)

type Status struct {
	Code    Code
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

func Equal(e1, e2 error) bool {
	return Error2Code(e1) == Error2Code(e2)
}

func Error2Code(err error) Code {
	// Don't use FromError to avoid allocation of OK status.
	if err == nil {
		return Ok
	}
	if se, ok := err.(*Status); ok {
		return se.Code
	}
	return Unknown
}

func Error2Status(err error) *Status {
	if err == nil {
		return &Status{
			Code:    Ok,
			Message: "Ok",
		}
	}

	if st, ok := err.(*Status); ok {
		return st
	}

	return &Status{
		Code:    Unknown,
		Message: err.Error(),
	}
}

type Code uint32

const (
	Ok                  Code = 0
	Unimplemented       Code = 1
	BadRequest          Code = 2
	NotFound            Code = 3
	PermissionDenied    Code = 4
	Internal            Code = 5
	Unauthenticated     Code = 6
	Unknown             Code = 7
	Unavailable         Code = 8
	InternalServerError Code = 9
	OutOfRange          Code = 10
	Max                 Code = 20
)

var code2Str = map[Code]string{
	Ok:                  "Ok",
	Unimplemented:       "Unimplemented",
	BadRequest:          "BadRequest",
	NotFound:            "NotFound",
	PermissionDenied:    "PermissionDenied",
	Internal:            "Internal",
	Unauthenticated:     "Unauthenticated",
	Unknown:             "Unknown",
	Unavailable:         "Unavailable",
	InternalServerError: "InternalServerError",
	OutOfRange:          "OutOfRange",
}

func (c Code) String() string {
	str, ok := code2Str[c]
	if !ok {
		return "Unknown"
	}

	return str
}

var (
	ErrUnimplemented       = New(Unimplemented, "Unimplemented")
	ErrBadRequest          = New(BadRequest, "BadRequest")
	ErrNotFound            = New(NotFound, "NotFound")
	ErrPermissionDenied    = New(PermissionDenied, "PermissionDenied")
	ErrInternal            = New(Internal, "Internal")
	ErrUnauthenticated     = New(Unauthenticated, "Unauthenticated")
	ErrUnknown             = New(Unknown, "Unknown")
	ErrUnavailable         = New(Unavailable, "Unavailable")
	ErrInternalServerError = New(InternalServerError, "InternalServerError")
	ErrOutOfRange          = New(OutOfRange, "OutOfRange")
)
