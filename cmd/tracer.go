package cmd

import (
	"fmt"
	"runtime"
)

// stackTraceEntry はスタックトレースの構造体
type StackTraceEntry struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Func string `json:"func"`
}

// StackTraceError はスタックトレースを含むエラー
type StackTraceError struct {
	Message    string            `json:"message"`
	StackTrace []StackTraceEntry `json:"stackTrace"`
}

// Error はエラーを文字列に変換する
func (e *StackTraceError) Error() string {
	starckTraceStr := ""
	for _, entry := range e.StackTrace {
		starckTraceStr += fmt.Sprintf("File: %s, Line: %d,  Func:%s\n", entry.File, entry.Line, entry.Func)
	}
	return fmt.Sprintf("%s\nStackTrace:\n%s", e.Message, starckTraceStr)
}

// 1. 新規のエラーを生成する NewStackTraceError("error message")
// 2. 既存のエラーにスタックトレースを追加する NewStackTraceError2(err, "")
//     errorが既にスタックトレースを含んでいる場合は追加しない message = err.Error()

// NewStackTraceError はスタックトレースを含むエラーを生成する
func NewStackTraceError(msg string) error {
	pc := make([]uintptr, 50)
	n := runtime.Callers(2, pc) // runtimeとNewStackTraceErrorをスキップ
	frames := runtime.CallersFrames(pc[:n])

	var stackTrace []StackTraceEntry
	for {
		frame, more := frames.Next()
		stackTrace = append(stackTrace, StackTraceEntry{
			File: frame.File,
			Line: frame.Line,
			Func: frame.Function,
		})
		if !more {
			break
		}
	}

	return &StackTraceError{
		Message:    msg,
		StackTrace: stackTrace,
	}
}

// WrapError はエラーにスタックトレースを追加する
// NewStackTraceError で生成したエラーには追加できない
func WrapError(err error) error {
	if err == nil {
		return nil
	}
	if stErr, ok := err.(*StackTraceError); ok {
		return stErr
	}
	pc := make([]uintptr, 50)
	n := runtime.Callers(2, pc) // runtimeとNewStackTraceErrorをスキップ
	frames := runtime.CallersFrames(pc[:n])

	var stackTrace []StackTraceEntry
	for {
		frame, more := frames.Next()
		stackTrace = append(stackTrace, StackTraceEntry{
			File: frame.File,
			Line: frame.Line,
			Func: frame.Function,
		})
		if !more {
			break
		}
	}

	return &StackTraceError{
		Message:    err.Error(),
		StackTrace: stackTrace,
	}
}
