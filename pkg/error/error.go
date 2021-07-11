package error

import "log"

var errors = map[int]string{}

type Error struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func NewError(code int, message string) *Error {
	if _, ok := errors[code]; ok {
		log.Fatalf("错误码 %d 已存在，请换一个", code)
	}

	errors[code] = message
	return &Error{Code: code, Message: message}
}

func (e *Error) WithDetails(details ...string) *Error {
	e.Details = append(e.Details, details...)
	return e
}
