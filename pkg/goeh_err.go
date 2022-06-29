package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Simon-Co/goeh/internal/calltrace"
)

var (
	ErrList = errors.New("Error List")
	ErrTest = errors.New("Test Error")
)

type GoehErr struct {
	File      string
	Operation string
	Line      int
	Message   string
	Cause     error
	trace     []string
	ErrorList *GoehErrList
}

// returns farmatted error string
func (e *GoehErr) Error() string {
	if errors.Is(e, ErrList) && e.ErrorList != nil {
		if len(e.ErrorList.List) > 0 {
			var els strings.Builder
			els.WriteString(fmt.Sprintf("File: %s\nOperation: %s\nMessage: %s\nError: %q\nTrace: [\n%s]\nErrorList:[", e.File, e.Operation, e.Message, e.Cause, e.trace))
			for _, ae := range e.ErrorList.List {
				els.WriteString(fmt.Sprintf("\nFile: %s\nOperation: %s\nMessage: %s\nError: %q\nTrace: [\n%s],\n", ae.File, ae.Operation, ae.Message, ae.Cause, ae.trace))
			}
			els.WriteString("]\n")
			return els.String()
		} else {
			return "No errors is ErrList"
		}
	}
	return fmt.Sprintf("File: %s\nOperation: %s\nLine: %d\nMessage: %s\nError: %q\nTrace: [\n%s]", e.File, e.Operation, e.Line, e.Message, e.Cause, e.trace)
}

//used in errors.Is check to see if error is of type Error
func (e *GoehErr) Is(target error) bool {
	t, ok := target.(*GoehErr)
	if !ok {
		return false
	}
	return t.Message == t.Message
}

//used in errors.Has
func (e *GoehErr) Unwrap() error {
	return e.Cause
}

//formats error trace string and adds to Error trace slice
func (e *GoehErr) addTrace(file string, operation string, line int) {
	trace := fmt.Sprintf("\nFile: %s; Operation: %s; Line: %d;\n", file, operation, line)
	e.trace = append(e.trace, trace)
}

// returns farmatted error string
func (e *GoehErr) String() string {
	if errors.Is(e, ErrList) && e.ErrorList != nil {
		if len(e.ErrorList.List) > 0 {
			var els strings.Builder
			els.WriteString(fmt.Sprintf("File: %s\nOperation: %s\nMessage: %s\nError: %q\nTrace: [\n%s]\nErrorList:[", e.File, e.Operation, e.Message, e.Cause, e.trace))
			for _, ae := range e.ErrorList.List {
				els.WriteString(fmt.Sprintf("\nFile: %s\nOperation: %s\nMessage: %s\nError: %q\nTrace: [\n%s],\n", ae.File, ae.Operation, ae.Message, ae.Cause, ae.trace))
			}
			els.WriteString("]\n")
			return els.String()
		} else {
			return "No errors is ErrList"
		}
	}
	return fmt.Sprintf("File: %s\nOperation: %s\nMessage: %s\nError: %q\nTrace: [\n%s]", e.File, e.Operation, e.Message, e.Cause, e.trace)
}

type GoehErrList struct {
	Mu   sync.Mutex
	List []*GoehErr
}

func (ael *GoehErrList) AddErr(e error) {
	ae := ParseToDepth(e, 3)
	ael.Mu.Lock()
	defer ael.Mu.Unlock()
	ael.List = append(ael.List, ae)
}

func NewErrorList() *GoehErr {
	return &GoehErr{Message: ErrList.Error(), Cause: ErrList, ErrorList: &GoehErrList{}}
}

//parses the supplied error.  If the error is of type *Error the
//error is added as a trace and is returned.  Else, the a new
// *Error is created using the suppiled error as a base.
func Parse(e error) *GoehErr {
	trace, _ := calltrace.Full(2)
	t, ok := e.(*GoehErr)
	if !ok {
		ae := &GoehErr{
			File:      trace.File,
			Operation: trace.Function,
			Line:      trace.Line,
			Message:   e.Error(),
			Cause:     e,
		}
		ae.addTrace(ae.File, ae.Operation, ae.Line)
		return ae
	} else {
		t.addTrace(trace.File, trace.Function, trace.Line)
		return t
	}
}

//parses the supplied error to the supplied callstaca depth.
//If the error is of type *Error the
//error is added as a trace and is returned.  Else, the a new
// *Error is created using the suppiled error as a base.
func ParseToDepth(e error, depth int) *GoehErr {
	trace, _ := calltrace.Full(depth)
	t, ok := e.(*GoehErr)
	if !ok {
		ae := &GoehErr{
			File:      trace.File,
			Operation: trace.Function,
			Line:      trace.Line,
			Message:   e.Error(),
			Cause:     e,
		}
		ae.addTrace(ae.File, ae.Operation, ae.Line)
		return ae
	} else {
		t.addTrace(trace.File, trace.Function, trace.Line)
		return t
	}
}
