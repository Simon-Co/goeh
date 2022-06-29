package calltrace

import (
	p "path/filepath"
	"runtime"
)

//returns the filepath relative to the supplied depth of the current
//runtime.Caller stack of where this function is called.
//depth of 1 will return the filepath relative to the initial caller of this function.
//depth of 0 will return the filepath relative to this function.
//higher depths will travel further back through the runtime.Caller stack,
//A failure to parse the runtime stack will return a ErrCannotPassRuntimeStack
func Filepath(depth int) (string, error) {
	_, filepath, _, ok := runtime.Caller(depth)
	if !ok {
		return "", ErrParseRuntimeCallerStack
	}
	return filepath, nil
}

//returns the Directory path relative to the supplied depth of the current
//runtime.Caller stack of where this function is called.
//depth of 1 will return relative to the initial caller of this function.
//depth of 0 will return the path relative to to the directory containing the Dirpath function.
//higher depths will travel further back through the runtime.Caller stack.
//A failure to parse the runtime stack will return a ErrCannotPassRuntimeStack
func Dirpath(depth int) (string, error) {
	_, filepath, _, ok := runtime.Caller(depth)
	if !ok {
		return "", ErrParseRuntimeCallerStack
	}
	return p.Dir(filepath), nil
}

//returns the Filename relative to the supplied depth of the current
//runtime.Caller stack of where this function is called.
//depth of 1 will return relative to the initial caller of this function.
//depth of 0 will return relative to the name of the file containing the filename function.
//higher depths will travel further back through the runtime.Caller stack.
//A failure to parse the runtime stack will return a ErrCannotPassRuntimeStack
func Filename(depth int) (string, error) {
	_, filepath, _, ok := runtime.Caller(depth)
	if !ok {
		return "", ErrParseRuntimeCallerStack
	}
	return p.Base(filepath), nil
}

//full runtime caller information
type FullTrace struct {
	File     string
	Dir      string
	Function string
	Line     int
}

//returns &FullTrace relative to the supplied depth of the current
//runtime.Caller stack of where this function is called.
//depth of 1 will return relative to the initial caller of this function.
//depth of 0 will return relative to the name of the file containing the filename function.
//higher depths will travel further back through the runtime.Caller stack.
//A failure to parse the runtime stack will return a ErrCannotPassRuntimeStack
func Full(depth int) (*FullTrace, error) {
	pc, fp, line, ok := runtime.Caller(depth)
	if !ok {
		return nil, ErrParseRuntimeCallerStack
	}
	callerFunc := runtime.FuncForPC(pc)
	return &FullTrace{
		File:     p.Base(fp),
		Dir:      fp,
		Function: p.Base(callerFunc.Name()),
		Line:     line,
	}, nil
}
