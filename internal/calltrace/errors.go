package calltrace

import "errors"

var (
	ErrParseRuntimeCallerStack = errors.New("unable to pass the runtime.Caller stack")
)
