package fastgo

type Error struct {
	error
	msg string
}

func NewError(str string) *Error {
	return &Error{msg: str}
}

func (e *Error) SetMsg(msg string) {
	e.msg = msg
}

func (e *Error) Error() string {
	return e.msg
}
