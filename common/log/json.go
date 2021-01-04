package log

import (
	"encoding/json"
	"fmt"
)

func Debugj(fmt string, args ...interface{}) {
	defaultLogger.Debugj(fmt, args...)
}
func Infoj(fmt string, args ...interface{}) {
	defaultLogger.Infoj(fmt, args...)
}
func Warnj(fmt string, args ...interface{}) {
	defaultLogger.Warnj(fmt, args...)
}
func Errorj(fmt string, args ...interface{}) {
	defaultLogger.Errorj(fmt, args...)
}

func toJson(args []interface{}) []interface{} {
	var newArg = make([]interface{}, 0, len(args))
	for i, arg := range args {
		bytes, err := json.Marshal(arg)
		if err != nil {
			newArg = append(newArg, fmt.Sprintf("[< arg %d = %#v | json.Marshal err = %+v >]", i, arg, err))
		} else {
			newArg = append(newArg, fmt.Sprintf("%s", bytes))
		}
	}
	return newArg
}
