package ethx

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"math/rand"
	"reflect"
	"regexp"
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

// AddressSlice parse any to []common.Address
func AddressSlice(addressLikeArr any) (addresses []common.Address) {
	arrValue := reflect.ValueOf(addressLikeArr)
	if arrValue.Kind() != reflect.Slice {
		panic(errors.New(fmt.Sprintf("param is not a slice: %v", addressLikeArr)))
	}
	n := arrValue.Len()
	addresses = make([]common.Address, n)
	for i := 0; i < n; i++ {
		addresses[i] = Address(arrValue.Index(i).Interface())
	}
	return addresses
}

// HashSlice parse any to []common.Hash
func HashSlice(hashLikeArr any) (hashes []common.Hash) {
	arrValue := reflect.ValueOf(hashLikeArr)
	if arrValue.Kind() != reflect.Slice {
		panic(errors.New(fmt.Sprintf("param is not a slice: %v", hashLikeArr)))
	}
	n := arrValue.Len()
	hashes = make([]common.Hash, n)
	for i := 0; i < n; i++ {
		hashes[i] = Hash(arrValue.Index(i).Interface())
	}
	return hashes
}

// BigIntSlice parse any to []*big.Int
func BigIntSlice(bigLikeArr any) (bigInts []*big.Int) {
	arrValue := reflect.ValueOf(bigLikeArr)
	if arrValue.Kind() != reflect.Slice {
		panic(errors.New(fmt.Sprintf("param is not a slice: %v", bigLikeArr)))
	}
	n := arrValue.Len()
	bigInts = make([]*big.Int, n)
	for i := 0; i < n; i++ {
		bigInts[i] = BigInt(arrValue.Index(i).Interface())
	}
	return bigInts
}

// TypeSlice assert any to []T
func TypeSlice[T any](arr any) []T {
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() != reflect.Slice {
		panic(errors.New(fmt.Sprintf("param is not a slice: %v", arr)))
	}
	n := arrValue.Len()
	vals := make([]T, n)
	for i := 0; i < n; i++ {
		vals[i] = arrValue.Index(i).Interface().(T)
	}
	return vals
}

var (
	r10, _ = regexp.Compile(`^[0-9]+$`)
	r2, _  = regexp.Compile(`^0[bB][01]+$`)
	r8, _  = regexp.Compile(`^0[oO][0-7]+$`)
	r16, _ = regexp.Compile(`^0[xX][0-9a-fA-F]+$`)

	bigInt10 = big.NewInt(10)
)

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
	case []byte:
		s := string(value)
		switch {
		case r10.MatchString(s):
			return stringBig(s, 10)
		case r16.MatchString(s):
			return stringBig(s[2:], 16)
		case r2.MatchString(s):
			return stringBig(s[2:], 2)
		case r8.MatchString(s):
			return stringBig(s[2:], 8)
		default:
			return new(big.Int).SetBytes(value)
		}
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
			return stringBig(value[2:], 16)
		case Is0b(value):
			return stringBig(value[2:], 2)
		case Is0o(value):
			return stringBig(value[2:], 8)
		default:
			for i := range value {
				switch value[i] {
				case 'e', 'E':
					return Mul(value[:i], new(big.Int).Exp(bigInt10, stringBig(value[i+1:], 10), nil))
				case '^':
					return Mul(new(big.Int).Exp(BigInt(value[:i]), BigInt(value[i+1:]), nil))
				}
			}
			return stringBig(value, 10)
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

func stringBig(value string, base int) *big.Int {
	v, b := new(big.Int).SetString(value, base)
	if !b {
		log.Panicf("Unknown numLike: %v", value)
	}
	return v
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

func MulDiv(numLike0, numLike1, numLike2 any) *big.Int {
	r := new(big.Int).Mul(BigInt(numLike0), BigInt(numLike1))
	return r.Div(r, BigInt(numLike2))
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
		return BigInt(numLike0).CmpAbs(BigInt(numLike1)) >= 0
	} else {
		return BigInt(numLike0).Cmp(BigInt(numLike1)) >= 0
	}
}

func Gt(numLike0, numLike1 any, isAbs ...bool) bool {
	if len(isAbs) > 0 && isAbs[0] {
		return BigInt(numLike0).CmpAbs(BigInt(numLike1)) == 1
	} else {
		return BigInt(numLike0).Cmp(BigInt(numLike1)) == 1
	}
}

func Lte(numLike0, numLike1 any, isAbs ...bool) bool {
	if len(isAbs) > 0 && isAbs[0] {
		return BigInt(numLike0).CmpAbs(BigInt(numLike1)) <= 0
	} else {
		return BigInt(numLike0).Cmp(BigInt(numLike1)) <= 0
	}
}

func Lt(numLike0, numLike1 any, isAbs ...bool) bool {
	if len(isAbs) > 0 && isAbs[0] {
		return BigInt(numLike0).CmpAbs(BigInt(numLike1)) == -1
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
