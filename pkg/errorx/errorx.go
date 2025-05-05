package errorx

import (
	"errors"
	"fmt"
	"reflect"
)

type ErrorCode uint64

const (
	CodeSuccess          ErrorCode = 0
	CodeInvalidArgument  ErrorCode = 40001
	CodePermissionDenied ErrorCode = 40002
	CodeInternalServer   ErrorCode = 50001
	CodeTooManyRequests  ErrorCode = 50002
)

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInternalServer   = errors.New("internal server error")
	ErrTooManyRequests  = errors.New("too many requests")
)

type Error struct {
	Err  error
	Code ErrorCode
	Meta interface{}
}

var _ error = (*Error)(nil)

func (msg *Error) SetCode(code ErrorCode) *Error {
	msg.Code = code
	return msg
}

// implements the error interface.
func (msg *Error) Error() string {
	return msg.Err.Error()
}

func (msg *Error) Unwrap() error {
	return msg.Err
}

func (msg *Error) SetMeta(data interface{}) *Error {
	msg.Meta = data
	return msg
}

func (msg *Error) JSON() interface{} {
	jsonData := make(map[string]interface{})
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				jsonData[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			jsonData["meta"] = msg.Meta
		}
	}
	if _, ok := jsonData["error"]; !ok {
		jsonData["error"] = msg.Error()
	}
	return jsonData
}

func New(err error, c ErrorCode, meta interface{}) *Error {
	return &Error{
		Err:  err,
		Code: c,
		Meta: meta,
	}
}

func NewArgumente(err error) *Error {
	return New(err, CodeInvalidArgument, nil)
}

func NewInternale(err error) *Error {
	return New(err, CodeInternalServer, nil)
}

func NewArgument(err string) *Error {
	return New(errors.New(err), CodeInvalidArgument, nil)
}

func NewInternal(err string) *Error {
	return New(errors.New(err), CodeInternalServer, nil)
}

func Newf(c ErrorCode, meta interface{}, format string, v ...interface{}) *Error {
	return New(fmt.Errorf(format, v...), c, meta)
}

func NewArgumentf(format string, v ...interface{}) *Error {
	return New(fmt.Errorf(format, v...), CodeInvalidArgument, nil)
}

func NewInternalf(format string, v ...interface{}) *Error {
	return New(fmt.Errorf(format, v...), CodeInternalServer, nil)
}
