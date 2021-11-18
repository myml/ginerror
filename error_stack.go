package ginerror

// ErrorStack 用于记录错误的堆栈信息
type ErrorStack struct {
	err      error
	function string
	filename string
}

func (es *ErrorStack) Unwrap() error {
	return es.err
}
func (es ErrorStack) Error() string {
	return es.err.Error()
}
