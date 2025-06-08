package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const (
	escape  = "\x1b["
	magenta = escape + "35m"
	def     = escape + "0m"
)

func MyErrorf(format string, a ...any) error {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf(format, a...)
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := filepath.Base(file)
	prefix := fmt.Sprintf("%s%s%s#%s%s%s", magenta, fileName, def, magenta, funcName, def)

	a = append([]any{prefix}, a...)

	return fmt.Errorf("%s "+format, a...)
}
