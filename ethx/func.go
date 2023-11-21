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
	methodName := getFuncName(f)
	funcValue := reflect.ValueOf(structInstance).MethodByName(methodName)
	if !funcValue.IsValid() {
		panic(errors.New(fmt.Sprintf("Method not found: %v", methodName)))
	}
	return call(funcValue, args...)
}

func getFuncName(f any) string {
	var methodName string
	switch f.(type) {
	case string:
		methodName = f.(string)
	default:
		methodName = runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	}
	if dotIndex := strings.LastIndex(methodName, "."); dotIndex != -1 {
		methodName = strings.Split(methodName[dotIndex+1:], "-")[0]
	}
	return methodName
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
