package ethx

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// callFunc call func with args
// structInstance could be nil
func callFunc(f any, args ...any) []any {
	return parseAny(call(reflect.ValueOf(f), args...))
}

func callStructMethod(structInstance any, f any, args ...any) []any {
	var funcValue reflect.Value
	funcName := getFuncName(f)
	switch v := structInstance.(type) {
	case reflect.Value:
		funcValue = v.MethodByName(funcName)
	default:
		funcValue = reflect.ValueOf(structInstance).MethodByName(funcName)
	}
	if !funcValue.IsValid() {
		panic(errors.New(fmt.Sprintf("Method not found: %v", funcName)))
	}
	return parseAny(call(funcValue, args...))
}

func call(funcValue reflect.Value, args ...any) []reflect.Value {
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
	return funcValue.Call(in)
}

func parseAny(ret []reflect.Value) []any {
	retArgs := make([]interface{}, len(ret))
	for i, v := range ret {
		retArgs[i] = v.Interface()
	}
	return retArgs
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

func getStructName(s any) (name string) {
	protoType := reflect.TypeOf(s)
	protoKind := protoType.Kind()
	switch protoKind {
	case reflect.Pointer:
		if protoKind == reflect.Pointer {
			protoType = protoType.Elem()
		}
		name = protoType.Name()
	case reflect.Struct:
		name = protoType.Name()
	case reflect.String:
		name = s.(string)
	default:
		panic(fmt.Errorf("%v::not support kind(%v)", getFuncName(getStructName), protoKind))
	}
	return name
}
