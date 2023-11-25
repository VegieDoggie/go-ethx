package ethx

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"math/rand"
	"reflect"
	"regexp"
	"time"
)

// Hash hashLike is non-nil
func Hash(hashLike any) common.Hash {
	if hashLike == nil {
		return common.Hash{}
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
func Address(addressLike any, isPri ...bool) common.Address {
	if addressLike == nil {
		return common.Address{}
	}
	if len(isPri) > 0 && isPri[0] {
		return crypto.PubkeyToAddress(PrivateKey(addressLike).PublicKey)
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

// AddressSlice parse any to []common.Address, eg: []string to []common.Address
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

// HashSlice parse any to []common.Hash, eg: []string to []common.Hash
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

// BigIntSlice parse any to []*big.Int, eg: []string to []*big.Int
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

// Type smart assert any to T
func Type[T any](x any) T {
	return abi.ConvertType(x, *new(T)).(T)
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
			return _stringBig(s, 10)
		case r16.MatchString(s):
			return _stringBig(s[2:], 16)
		case r2.MatchString(s):
			return _stringBig(s[2:], 2)
		case r8.MatchString(s):
			return _stringBig(s[2:], 8)
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
		return stringBig(value)
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

func stringBig(value string) *big.Int {
	switch {
	case Is0x(value):
		return _stringBig(value[2:], 16)
	case Is0b(value):
		return _stringBig(value[2:], 2)
	case Is0o(value):
		return _stringBig(value[2:], 8)
	default:
		for i := range value {
			switch value[i] {
			case 'e', 'E':
				return Mul(value[:i], new(big.Int).Exp(bigInt10, _stringBig(value[i+1:], 10), nil))
			}
		}
		return _stringBig(value, 10)
	}
}

func _stringBig(value string, base int) *big.Int {
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

var rpcRegx, _ = regexp.Compile(`((?:https|wss|http|ws)[^\s\n\\"]+)`)

// CheckRpcLogged returns what rpcs are reliable for filter logs
// example:
//  1. rpc list: CheckRpcLogged("https://bsc-dataseed1.defibit.io", "https://bsc-dataseed4.binance.org")
//  2. auto resolve rpc list: CheckRpcLogged("https://bsc-dataseed1.defibit.io\t29599361\t1.263s\t\t\nConnect Wallet\nhttps://bsc-dataseed4.binance.org")
func CheckRpcLogged(rpcLike ...string) (reliableList []string) {
	var query = ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(0),
		ToBlock:   new(big.Int).SetUint64(1000),
	}
	log.Println("CheckRpcLogged start......")
	for _, iRpc := range rpcLike {
		rpcList := rpcRegx.FindAllString(iRpc, -1)
		for _, jRpc := range rpcList {
			client, err := ethclient.Dial(jRpc)
			if err == nil {
				logs, err := client.FilterLogs(context.TODO(), query)
				if err == nil && len(logs) > 0 {
					reliableList = append(reliableList, jRpc)
					continue
				}
			}
			log.Printf("[WARN] CheckRpcLogged::Unreliable rpc: %v\n", jRpc)
		}
	}
	log.Println("CheckRpcLogged finished......")
	return reliableList
}

// CheckRpcSpeed returns the rpc speed list
// example:
//  1. CheckRpcSpeed("https://bsc-dataseed1.defibit.io", "https://bsc-dataseed4.binance.org")
//  2. CheckRpcSpeed("https://bsc-dataseed1.defibit.io\t29599361\t1.263s\t\t\nConnect Wallet\nhttps://bsc-dataseed4.binance.org")
func CheckRpcSpeed(rpcLike ...string) (rpcSpeedMap map[string]time.Duration) {
	rpcSpeedMap = make(map[string]time.Duration)
	log.Println("CheckRpcSpeed start......")
	for _, iRpc := range rpcLike {
		rpcList := rpcRegx.FindAllString(iRpc, -1)
		for _, jRpc := range rpcList {
			client, err := ethclient.Dial(jRpc)
			if err == nil {
				before := time.Now()
				blockNumber, err := client.BlockNumber(context.TODO())
				if err == nil && blockNumber > 0 {
					rpcSpeedMap[jRpc] = time.Since(before)
					log.Printf("[%v] %v\r\n", rpcSpeedMap[jRpc], jRpc)
					continue
				}
			}
		}
	}
	log.Println("CheckRpcSpeed finished......")
	return rpcSpeedMap
}

func PrivateKey(priLike any) *ecdsa.PrivateKey {
	var privateStr string
	switch value := priLike.(type) {
	case *ecdsa.PrivateKey:
		return value
	case string:
		if Is0x(value) {
			privateStr = value[2:]
		} else {
			privateStr = value
		}
	default:
		privateStr = common.BigToHash(BigInt(value)).Hex()[2:]
	}
	privateKey, err := crypto.HexToECDSA(privateStr)
	if err != nil {
		panic(err)
	}
	return privateKey
}
