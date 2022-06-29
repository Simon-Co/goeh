package calltrace

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetFilePath(t *testing.T) {
	_, filepath, _, _ := runtime.Caller(0)
	tc := struct {
		expected string
		error    error
	}{
		expected: filepath,
		error:    nil,
	}
	result, err := Filepath(1)
	if err != nil {
		if tc.error == nil {
			t.Fatalf("\nUnexpected Error\nExpected: %s\nReceived: %s", tc.error.Error(), err.Error())
		}
	}
	if result != tc.expected {
		t.Fatalf("\nExpected Result\nExpected: %s\nReceived: %s", tc.expected, result)
	}
}

func TestGetDirPath(t *testing.T) {
	_, fp, _, _ := runtime.Caller(0)
	tc := struct {
		expected string
		error    error
	}{
		expected: filepath.Dir(fp),
		error:    nil,
	}
	result, err := Dirpath(1)
	if err != nil {
		if tc.error == nil {
			t.Fatalf("\nUnexpected Error\nExpected: %s\nReceived: %s", tc.error.Error(), err.Error())
		}
	}
	if result != tc.expected {
		t.Fatalf("\nExpected Result\nExpected: %s\nReceived: %s", tc.expected, result)
	}
}

func TestGetFilename(t *testing.T) {
	_, fp, _, _ := runtime.Caller(0)
	tc := struct {
		expected string
		error    error
	}{
		expected: filepath.Base(fp),
		error:    nil,
	}
	result, err := Filename(1)
	if err != nil {
		if tc.error == nil {
			t.Fatalf("\nUnexpected Error\nExpected: %s\nReceived: %s", tc.error.Error(), err.Error())
		}
	}
	if result != tc.expected {
		t.Fatalf("\nExpected Result\nExpected: %s\nReceived: %s", tc.expected, result)
	}
}

func TestGetFull(t *testing.T) {
	_, fp, _, _ := runtime.Caller(0)
	testCases := []struct {
		name     string
		input    int
		expected FullTrace
		error    error
	}{
		{
			name:  "ParseFail",
			input: 5,
			error: ErrParseRuntimeCallerStack,
		},
		{
			name:  "Success",
			input: 1,
			expected: FullTrace{
				File: filepath.Base(fp),
				Dir:  fp,
				Line: 98,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pc, _, _, _ := runtime.Caller(0)
			callerFunc := runtime.FuncForPC(pc)
			result, err := Full(tc.input)
			if err != nil {
				if tc.error == nil || !errors.Is(err, tc.error) {
					t.Fatalf("\nUnexpected Error\nExpected: %s\nReceived: %s", tc.error.Error(), err.Error())
				}
				return
			}
			tc.expected.Function = filepath.Base(callerFunc.Name())
			if tc.expected != *result {
				t.Errorf("\nExpected Result\nExpected: %+v\nReceived: %+v", tc.expected, *result)
			}
		})
	}
}
