package error

import "fmt"

type Error struct {
	code    int		`json:"code"`
	msg     string	`json:"msg"`
	details []string `json:"details"`
}

var codes = map[int]string{} 

func NewError(code int, msg string) *Error {
	if v, ok := codes[code]; ok {
		panic(fmt.Sprintf("code: %d is allready existed msg%s\n", code, v))
	}	

	codes[code] = msg

	return &Error{
		code: code,
		msg: msg,
		details: make([]string, 0),
	}
}


func (e *Error)Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code(), e.Msg())
}

func (e *Error)Code() int {
	return e.code
}

func (e *Error)Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	for _, d := range details {
		newError.details = append(newError.details, d)
	}

	return &newError
}

