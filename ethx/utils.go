package ethx

import (
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"math/rand"
	"reflect"
)

// Hash hashLike is non-nil
func Hash(hashLike any) common.Hash {
	if hashLike == nil {
		log.Panicf("hashLike is non-nil")
	}
	switch value := hashLike.(type) {
	case common.Hash:
		return value
	case *common.Hash:
		return *value
	default:
		return common.BigToHash(BigInt(hashLike))
	}
}

// Address addressLike is non-nil
func Address(addressLike any) common.Address {
	if addressLike == nil {
		log.Panicf("addressLike is non-nil")
	}
	switch value := addressLike.(type) {
	case common.Address:
		return value
	case *common.Address:
		return *value
	default:
		return common.BigToAddress(BigInt(addressLike))
	}
}

func BigInt(numLike any) *big.Int {
	if numLike == nil {
		return nil
	}
	switch value := numLike.(type) {
	case *big.Int:
		return value
	case *common.Address:
		return value.Big()
	case common.Address:
		return value.Big()
	case *common.Hash:
		return value.Big()
	case common.Hash:
		return value.Big()
	case *[]byte:
		return new(big.Int).SetBytes(*value)
	case []byte:
		return new(big.Int).SetBytes(value)
	case int:
		return big.NewInt(int64(value))
	case int8:
		return big.NewInt(int64(value))
	case int16:
		return big.NewInt(int64(value))
	case int32:
		return big.NewInt(int64(value))
	case int64:
		return big.NewInt(value)
	case uint:
		return new(big.Int).SetUint64(uint64(value))
	case uint8:
		return new(big.Int).SetUint64(uint64(value))
	case uint16:
		return new(big.Int).SetUint64(uint64(value))
	case uint32:
		return new(big.Int).SetUint64(uint64(value))
	case uint64:
		return new(big.Int).SetUint64(value)
	case float32:
		return new(big.Int).SetUint64(uint64(value))
	case float64:
		return new(big.Int).SetUint64(uint64(value))
	case string:
		switch {
		case Is0x(value):
			v, b := new(big.Int).SetString(value[2:], 16)
			if !b {
				log.Panicf("Unknown numLike: %v", value)
			}
			return v
		case Is0b(value):
			v, b := new(big.Int).SetString(value[2:], 2)
			if !b {
				log.Panicf("Unknown numLike: %v", value)
			}
			return v
		case Is0o(value):
			v, b := new(big.Int).SetString(value[2:], 8)
			if !b {
				log.Panicf("Unknown numLike: %v", value)
			}
			return v
		default:
			v, b := new(big.Int).SetString(value, 10)
			if !b {
				log.Panicf("Unknown numLike: %v", value)
			}
			return v
		}
	case bool:
		if value {
			return big.NewInt(1)
		}
		return big.NewInt(0)
	default:
		rv := reflect.ValueOf(value)
		kind := rv.Kind()
		switch kind {
		case reflect.Pointer, reflect.UnsafePointer:
			return BigInt(rv.Elem().Interface())
		}
		log.Panicf("Unknown numLike: %v", value)
		return nil
	}
}

func Uint64(numLike any) uint64 {
	return BigInt(numLike).Uint64()
}

func Int64(numLike any) int64 {
	return BigInt(numLike).Int64()
}

func Add(numLike ...any) *big.Int {
	return Sum(numLike...)
}

func Sub(numBase any, numLike ...any) *big.Int {
	num := BigInt(numBase)
	for i := range numLike {
		num.Sub(num, BigInt(numLike[i]))
	}
	return num
}

func Mul(numLike ...any) *big.Int {
	num := big.NewInt(1)
	for i := range numLike {
		num.Mul(num, BigInt(numLike[i]))
	}
	return num
}

func Div(numLike0, numLike1 any) *big.Int {
	return new(big.Int).Div(BigInt(numLike0), BigInt(numLike1))
}

func Sum(numLike ...any) *big.Int {
	num := new(big.Int)
	for i := range numLike {
		num.Add(num, BigInt(numLike[i]))
	}
	return num
}

func Gte(numLike0, numLike1 any, isAbs ...bool) bool {
	if len(isAbs) > 0 && isAbs[0] {
		return BigInt(numLike0).CmpAbs(BigInt(numLike1)) <= 0
	} else {
		return BigInt(numLike0).Cmp(BigInt(numLike1)) <= 0
	}
}

func Gt(numLike0, numLike1 any, isAbs ...bool) bool {
	if len(isAbs) > 0 && isAbs[0] {
		return BigInt(numLike0).CmpAbs(BigInt(numLike1)) == -1
	} else {
		return BigInt(numLike0).Cmp(BigInt(numLike1)) == -1
	}
}

func Lte(numLike0, numLike1 any, isAbs ...bool) bool {
	if len(isAbs) > 0 && isAbs[0] {
		return BigInt(numLike0).CmpAbs(BigInt(numLike1)) >= 0
	} else {
		return BigInt(numLike0).Cmp(BigInt(numLike1)) == -1
	}
}

func Lt(numLike0, numLike1 any, isAbs ...bool) bool {
	if len(isAbs) > 0 && isAbs[0] {
		return BigInt(numLike0).CmpAbs(BigInt(numLike1)) == 1
	} else {
		return BigInt(numLike0).Cmp(BigInt(numLike1)) == -1
	}
}

func Is0x(s string) bool {
	return len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

func Is0b(s string) bool {
	return len(s) > 2 && s[0] == '0' && (s[1] == 'b' || s[1] == 'B')
}

func Is0o(s string) bool {
	return len(s) > 2 && s[0] == '0' && (s[1] == 'o' || s[1] == 'O')
}

func RandBytes(len int) []byte {
	b := make([]byte, len)
	for i := 0; i < len; i++ {
		b[i] = byte(rand.Int())
	}
	return b
}
