package ethx

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"testing"
)

var rawRPC = "https://api.zan.top/node/v1/bsc/testnet/public\t35403779\t0.077s\t\t\nAdd to Metamask\nhttps://data-seed-prebsc-1-s1.bnbchain.org:8545\t35403779\t0.087s\t\t\nAdd to Metamask\nhttps://data-seed-prebsc-1-s2.bnbchain.org:8545\t35403779\t0.088s\t\t\nAdd to Metamask\nhttps://bsc-testnet.blockpi.network/v1/rpc/public\t35403779\t0.090s\t\t\nAdd to Metamask\nhttps://data-seed-prebsc-2-s1.bnbchain.org:8545\t35403779\t0.100s\t\t\nAdd to Metamask\nhttps://data-seed-prebsc-2-s2.bnbchain.org:8545\t35403779\t0.100s\t\t\nAdd to Metamask\nhttps://endpoints.omniatech.io/v1/bsc/testnet/public\t35403779\t0.209s\t\t\nAdd to Metamask\nwss://bsc-testnet.publicnode.com\t35403779\t0.213s\t\t\nhttps://bsc-testnet.publicnode.com\t35403779\t0.243s\t\t\nAdd to Metamask\nhttps://bsc-testnet.public.blastapi.io\t35403779\t0.250s\t\t\nAdd to Metamask\nhttps://data-seed-prebsc-2-s3.bnbchain.org:8545\t\t\t\t\nAdd to Metamask\nhttps://data-seed-prebsc-1-s3.bnbchain.org:8545\t\t\t\t\nAdd to Metamask\nhttps://bsctestapi.terminet.io/rpc"

func Test_resolveRPCs(t *testing.T) {
	l1 := resolveRPCs(rawRPC)
	l2 := resolveRPCs(rawRPC, rawRPC, rawRPC)
	assert.Equal(t, len(l1), len(l2))
}

func Test_Hash(t *testing.T) {
	expect := common.HexToHash("0x1")
	decode, err := hex.DecodeString("01")
	if err != nil {
		panic(err)
	}
	data := []any{
		nil,
		expect,
		&expect,
		expect.String(),
		big.NewInt(1),
		[]byte("0x1"),
		[]byte("1"),
		[]byte("0b1"),
		[]byte("0o1"),
		decode,
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		float32(1),
		float64(1),
		"1",
		"1e0",
		true,
		&decode,
		&expect,
	}
	for i, tdd := range data {
		if tdd == nil {
			assert.Equal(t, common.Hash{}, Hash(tdd), fmt.Sprintf("[%v] expect=%v, tdd=%v", i, expect, tdd))
		} else {
			assert.Equal(t, expect, Hash(tdd), fmt.Sprintf("[%v] expect=%v, tdd=%v", i, expect, tdd))
		}
	}
}

func Test_Address(t *testing.T) {
	expect := common.HexToAddress("0x1")
	decode, err := hex.DecodeString("01")
	if err != nil {
		panic(err)
	}
	data := []any{
		nil,
		expect,
		&expect,
		expect.String(),
		big.NewInt(1),
		[]byte("0x1"),
		[]byte("1"),
		[]byte("0b1"),
		[]byte("0o1"),
		decode,
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		float32(1),
		float64(1),
		"1",
		"1e0",
		true,
		&decode,
		&expect,
	}
	for i, td := range data {
		if td == nil {
			assert.Equal(t, common.Address{}, Address(td), fmt.Sprintf("[%v] expect=%v, td=%v", i, expect, td))
		} else {
			assert.Equal(t, expect, Address(td), fmt.Sprintf("[%v] expect=%v, td=%v", i, expect, td))
		}
	}
	privateKey, err := crypto.HexToECDSA("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		panic(err)
	}
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)
	assert.Equal(t, addr, Address("0x0000000000000000000000000000000000000000000000000000000000000001", true))
	assert.Equal(t, addr, Address("0000000000000000000000000000000000000000000000000000000000000001", true))
	assert.Equal(t, addr, Address("1", true))
	assert.Equal(t, addr, Address(1, true))
}

func Test_AddressSlice(t *testing.T) {
	addrArr := []string{"1", "2", "3", "4"}
	addresses := AddressSlice(addrArr)
	for i := range addresses {
		assert.Equal(t, Address(addrArr[i]), addresses[i])
	}
}

func Test_HashSlice(t *testing.T) {
	addrArr := []string{"1", "2", "3", "4"}
	hashes := HashSlice(addrArr)
	for i := range hashes {
		assert.Equal(t, Hash(addrArr[i]), hashes[i])
	}
}

func Test_BigIntSlice(t *testing.T) {
	addrArr := []string{"1", "2", "3", "4"}
	bigInts := BigIntSlice(addrArr)
	for i := range bigInts {
		assert.Equal(t, BigInt(addrArr[i]).Uint64(), bigInts[i].Uint64())
	}
}

func Test_Uint64(t *testing.T) {
	addrArr := []string{"1", "2", "3", "4"}
	for i := range addrArr {
		parseUint, _ := strconv.ParseUint(addrArr[i], 10, 64)
		assert.Equal(t, Uint64(addrArr[i]), parseUint)
	}
}

func Test_Int64(t *testing.T) {
	addrArr := []string{"1", "2", "3", "4"}
	for i := range addrArr {
		parseUint, _ := strconv.ParseInt(addrArr[i], 10, 64)
		assert.Equal(t, Int64(addrArr[i]), parseUint)
	}
}

func Test_Add(t *testing.T) {
	for i := 1; i < 10; i++ {
		var params []int
		var sum int64
		for j := 0; j < i; j++ {
			params = append(params, 100+j)
			sum += int64(100 + j)
		}
		assert.Equal(t, sum, Add(params).Int64())
	}
	assert.Equal(t, int64(100+21), Add([]int{100}, 21).Int64())
	assert.Equal(t, int64(33+100+21), Add([]int{33, 100}, 21).Int64())
}

func Test_Mul(t *testing.T) {
	for i := 1; i < 10; i++ {
		var params []int
		mul := int64(1)
		for j := 0; j < i; j++ {
			params = append(params, 100+j)
			mul *= int64(100 + j)
		}
		assert.Equal(t, mul, Mul(params).Int64())
	}
	assert.Equal(t, int64(100*21), Mul([]int{100}, 21).Int64())
	assert.Equal(t, int64(33*100*21), Mul([]int{33, 100}, 21).Int64())
}

func Test_Div(t *testing.T) {
	assert.Equal(t, int64(100/35), Div(100, 35).Int64())
	assert.Equal(t, int64(100/101), Div(100, 101).Int64())
}

func Test_MulDiv(t *testing.T) {
	assert.Equal(t, int64(100*3/35), MulDiv(100, 3, 35).Int64())
	assert.Equal(t, int64(100*0/101), MulDiv(100, 0, 101).Int64())
}

func Test_Sub(t *testing.T) {
	assert.Equal(t, int64(1-100-3-35), Sub(1, 100, 3, 35).Int64())
	assert.Equal(t, int64(1200-100-0-101), Sub(1200, 100, 0, 101).Int64())
	assert.Equal(t, int64(1200-100-0-101), Sub(1200, []int{100, 0, 101}).Int64())
}

func Test_Gte(t *testing.T) {
	assert.Equal(t, false, Gte(100, 200))
	assert.Equal(t, true, Gte(200, 200))
	assert.Equal(t, true, Gte(300, 200))

	assert.Equal(t, false, Gte(-100, -200, true))
	assert.Equal(t, true, Gte(-200, -200, true))
	assert.Equal(t, true, Gte(-300, -200, true))
}

func Test_Gt(t *testing.T) {
	assert.Equal(t, false, Gt(100, 200))
	assert.Equal(t, false, Gt(200, 200))
	assert.Equal(t, true, Gt(300, 200))

	assert.Equal(t, false, Gt(-100, -200, true))
	assert.Equal(t, false, Gt(-200, -200, true))
	assert.Equal(t, true, Gt(-300, -200, true))
}

func Test_Lte(t *testing.T) {
	assert.Equal(t, true, Lte(100, 200))
	assert.Equal(t, true, Lte(200, 200))
	assert.Equal(t, false, Lte(300, 200))

	assert.Equal(t, true, Lte(-100, -200, true))
	assert.Equal(t, true, Lte(-200, -200, true))
	assert.Equal(t, false, Lte(-300, -200, true))
}

func Test_Lt(t *testing.T) {
	assert.Equal(t, true, Lt(100, 200))
	assert.Equal(t, false, Lt(200, 200))
	assert.Equal(t, false, Lt(300, 200))

	assert.Equal(t, true, Lt(-100, -200, true))
	assert.Equal(t, false, Lt(-200, -200, true))
	assert.Equal(t, false, Lt(-300, -200, true))
}

func Test_RandBytes(t *testing.T) {
	assert.Equal(t, 10*2, len(hex.EncodeToString(RandBytes(10))))
}

func Test_CheckRpcSpeed(t *testing.T) {
	//CheckRpcSpeed(rawRPC)
}

func Test_CheckRpcLogged(t *testing.T) {
	//CheckRpcLogged(rawRPC)
}
