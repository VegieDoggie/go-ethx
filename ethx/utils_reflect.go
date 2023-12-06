package ethx

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"reflect"
	"runtime"
	"strings"
)

// callFunc call func with args
// structInstance could be nil
func callFunc(f any, args ...any) []any {
	return parseAny(call(reflect.ValueOf(f), args))
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
	return parseAny(call(funcValue, args))
}

var (
	callOptsPtrType = reflect.TypeOf(new(bind.CallOpts))
	bigIntPtrType   = reflect.TypeOf(new(big.Int))
	addressType     = reflect.TypeOf(common.Address{}) // include [20]byte
	addressPtrType  = reflect.TypeOf(new(common.Address))
	hashType        = reflect.TypeOf(common.Hash{}) // include [32]byte
	hashPtrType     = reflect.TypeOf(new(common.Hash))
)

func call(funcValue reflect.Value, args []any) []reflect.Value {
	funcType := funcValue.Type()
	paramNum := funcType.NumIn()
	if paramNum != len(args) {
		panic(fmt.Errorf("call::Params num mismatch! Required: %v, but %v\n", paramNum, len(args)))
	}
	var in []reflect.Value
	for i := 0; i < paramNum; i++ {
		if inType := funcType.In(i); args[i] == nil {
			in = append(in, reflect.Zero(inType))
		} else {
			if argType := reflect.TypeOf(args[i]); !argType.ConvertibleTo(inType) {
				switch {
				case bigIntPtrType.ConvertibleTo(inType):
					args[i] = BigInt(args[i])
				case addressType.ConvertibleTo(inType):
					args[i] = Address(args[i])
				case addressPtrType.ConvertibleTo(inType):
					args[i] = AddressPtr(args[i])
				case hashType.ConvertibleTo(inType):
					args[i] = Hash(args[i])
				case hashPtrType.ConvertibleTo(inType):
					args[i] = HashPtr(args[i])
				default:
					panic(fmt.Errorf("[WARN] call(%v)::Param type mismatch! Required: %v, but %v\n", getFuncName(funcValue), inType, reflect.TypeOf(args[i])))
				}
			}
			in = append(in, reflect.ValueOf(args[i]))
		}
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
	case reflect.Value:
		methodName = runtime.FuncForPC(f.(reflect.Value).Pointer()).Name()
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
