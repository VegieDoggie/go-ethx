package ethx

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func callFunc(f any, args ...any) []any {
	return call(reflect.ValueOf(f), args...)
}

func callStructFunc(structInstance any, f any, args ...any) []any {
	var methodName string
	switch f.(type) {
	case string:
		methodName = f.(string)
	default:
		methodName = getFuncName(f)
	}
	funcValue := reflect.ValueOf(structInstance).MethodByName(methodName)
	if !funcValue.IsValid() {
		panic(errors.New(fmt.Sprintf("Method not found: %v", methodName)))
	}
	return call(funcValue, args...)
}

func getFuncName(f any) string {
	funcPtr := reflect.ValueOf(f).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	dotIndex := strings.LastIndex(funcName, ".")
	if dotIndex != -1 {
		funcName = strings.Split(funcName[dotIndex+1:], "-")[0]
	}
	return funcName
}

func call(funcValue reflect.Value, args ...any) []any {
	funcType := funcValue.Type()
	var in []reflect.Value
	for i, arg := range args {
		if arg == nil {
			in = append(in, reflect.Zero(funcType.In(i)))
		} else {
			in = append(in, reflect.ValueOf(arg))
		}
	}
	if len(in) == 0 && funcType.NumIn() == 1 {
		in = append(in, reflect.Zero(funcType.In(0)))
	}
	ret := funcValue.Call(in)
	retArgs := make([]interface{}, len(ret))
	for i, v := range ret {
		retArgs[i] = v.Interface()
	}
	return retArgs
}
