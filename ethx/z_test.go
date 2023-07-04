package ethx

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func TestBasicUsage(t *testing.T) {
	rpcList := []string{
		"https://data-seed-prebsc-2-s2.binance.org:8545",
		"https://data-seed-prebsc-2-s1.binance.org:8545",
	}
	weights := []int{1, 2}
	clientx := NewClientx(rpcList, weights, 30)

	// auto block number test
	go func() {
		queryTicker := time.NewTicker(time.Second)
		defer queryTicker.Stop()
		for {
			fmt.Println("AutoBlockNumber", clientx.AutoBlockNumber)
			<-queryTicker.C
		}
	}()
	select {}
}

func TestScanner(t *testing.T) {
	rpcList := []string{
		"https://data-seed-prebsc-2-s2.binance.org:8545",
		"https://data-seed-prebsc-2-s1.binance.org:8545",
	}
	clientx := NewClientx(rpcList, nil, 30)
	scanner := clientx.NewScanner(nil, nil, 200, 800, 3)
	logs, _ := scanner.Scan(0, 3000)
	fmt.Println(logs)
}

func TestLimit(t *testing.T) {
	defaultLimiter := rate.NewLimiter(5, 5)
	ctx := context.TODO()
	for i := 0; ; i++ {
		_ = defaultLimiter.Wait(ctx)
		fmt.Println(i)
	}
}

func TestUtils(t *testing.T) {
	// To Hash
	fmt.Println(Hash("1"))   // 0x0000000000000000000000000000000000000000000000000000000000000001
	fmt.Println(Hash(1))     // 0x0000000000000000000000000000000000000000000000000000000000000001
	fmt.Println(Hash("0x1")) // 0x0000000000000000000000000000000000000000000000000000000000000001
	fmt.Println(Hash("0b1")) // 0x0000000000000000000000000000000000000000000000000000000000000001
	fmt.Println(Hash("0o1")) // 0x0000000000000000000000000000000000000000000000000000000000000001

	// To Address
	fmt.Println(Address("1"))   // 0x0000000000000000000000000000000000000001
	fmt.Println(Address(1))     // 0x0000000000000000000000000000000000000001
	fmt.Println(Address("0x1")) // 0x0000000000000000000000000000000000000001
	fmt.Println(Address("0b1")) // 0x0000000000000000000000000000000000000001
	fmt.Println(Address("0o1")) // 0x0000000000000000000000000000000000000001

	// To BigInt
	fmt.Println(BigInt("1"))   // 1
	fmt.Println(BigInt(1))     // 1
	fmt.Println(BigInt("0x1")) // 1
	fmt.Println(BigInt("0b1")) // 1
	fmt.Println(BigInt("0o1")) // 1

	// + - * / 和 比较
	fmt.Println(Add("1", 2, 3, Hash(1)))     // 7
	fmt.Println(Sum("1", 2, 3, Hash(1)))     // 7
	fmt.Println(Sub(Address(1), Hash(0)))    // 1
	fmt.Println(Mul("0x2", 5))               // 10
	fmt.Println(Div("10", 2))                // 5
	fmt.Println(Gte(Address(1), Address(1))) // true
	fmt.Println(Gt(Address(1), Address(1)))  // false
	fmt.Println(Lte(Address(1), Address(1))) // true
	fmt.Println(Lt(Address(1), Address(1)))  // false
	fmt.Println(Uint64(Address(1)))          // 1
	fmt.Println(Int64(Address(1)))           // 1

	// 随机
	fmt.Println(Hash(RandBytes(32)))    // 0x06e4ba3da81342545e60108a576ef5590ee56800ef285bd692923b696f05fa44
	fmt.Println(Address(RandBytes(20))) // 0x8807B00e663fD283ff7e9C1291EFF9D6963290Da

	// Rpc 检查连通性和响应速率
	raw := "https://bsc-dataseed1.defibit.io\t29599361\t1.263s\t\t\nConnect Wallet\nhttps://bsc-dataseed4.binance.org\t29599361\t1.263s\t\t\nConnect Wallet\nhttps://bsc.blockpi.network/v1/rpc/public\t29599361\t1.263s\t\t\nConnect Wallet\nhttps://bsc-dataseed2.defibit.io\t29599361\t1.264s\t\t\nConnect Wallet\nhttps://bsc-dataseed1.binance.org\t29599361\t1.264s\t\t\nConnect Wallet\nhttps://koge-rpc-bsc.48.club\t29599361\t1.264s\t\t\nConnect Wallet\nhttps://bsc-dataseed3.binance.org\t29599361\t1.265s\t\t\nConnect Wallet\nhttps://bsc-dataseed4.defibit.io\t29599361\t1.548s\t\t\nConnect Wallet\nhttps://bsc-dataseed2.binance.org\t29599361\t1.548s\t\t\nConnect Wallet\nhttps://bsc-mainnet.rpcfast.com?api_key=xbhWBI1Wkguk8SNMu1bvvLurPGLXmgwYeC4S6g2H7WdwFigZSmPWVZRxrskEQwIf\t29599361\t1.548s\t\t\nConnect Wallet\nhttps://bsc-dataseed3.ninicoin.io\t29599361\t1.549s\t\t\nConnect Wallet\nhttps://bsc-mainnet.gateway.pokt.network/v1/lb/6136201a7bad1500343e248d\t29599361\t1.549s\t\t\nConnect Wallet\nhttps://bsc-dataseed2.ninicoin.io\t29599361\t1.840s\t\t\nConnect Wallet\nhttps://rpc-bsc.48.club\t29599361\t1.843s\t\t\nConnect Wallet\nhttps://rpc.ankr.com/bsc\t29599361\t1.843s\t\t\nConnect Wallet\nhttps://1rpc.io/bnb\t29599361\t1.843s\t\t\nConnect Wallet\nhttps://bsc-dataseed3.defibit.io\t29599361\t1.844s\t\t\nConnect Wallet\nhttps://bsc.publicnode.com\t29599361\t1.844s\t\t\nConnect Wallet\nhttps://bsc-dataseed4.ninicoin.io\t29599361\t2.115s\t\t\nConnect Wallet\nhttps://endpoints.omniatech.io/v1/bsc/mainnet/public\t29599361\t2.115s\t\t\nConnect Wallet\nhttps://bsc-mainnet.nodereal.io/v1/64a9df0874fb4a93b9d0a3849de012d3\t29599361\t2.115s\t\t\nConnect Wallet\nhttps://binance.nodereal.io\t29599361\t2.387s\t\t\nConnect Wallet\nhttps://bsc-mainnet.public.blastapi.io\t29599361\t2.387s\t\t\nConnect Wallet\nhttps://bsc.meowrpc.com\t29599361\t2.654s\t\t\nConnect Wallet\nhttps://bscrpc.com\t29599360\t1.549s\t\t\nConnect Wallet\nhttps://bnb.api.onfinality.io/public\t29599359\t1.845s\t\t\nConnect Wallet\nhttps://bsc-dataseed1.ninicoin.io\t29598929\t1.814s\t\t\nConnect Wallet\nwss://bsc-ws-node.nariox.org\t\t\t\t\nhttps://nodes.vefinetwork.org/smartchain\t\t\t\t\nConnect Wallet\nhttps://bsc.rpcgator.com"
	reliableRpcList, _, _, _ := CheckRpcConn(raw)
	CheckRpcSpeed(reliableRpcList...)
}
