package ethx

import (
	"fmt"
	"log"
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
	//defaultLimiter := rate.NewLimiter(5, 5)
	//ctx := context.TODO()
	//for i := 0; ; i++ {
	//	_ = defaultLimiter.Wait(ctx)
	//	fmt.Println(i)
	//}
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
	fmt.Println(BigInt("1"))     // 1
	fmt.Println(BigInt(1))       // 1
	fmt.Println(BigInt("0x1"))   // 1
	fmt.Println(BigInt("0b1"))   // 1
	fmt.Println(BigInt("0o1"))   // 1
	fmt.Println(BigInt("1e3"))   // 1000
	fmt.Println(BigInt("10E3"))  // 10000
	fmt.Println(BigInt("2^3"))   // 8
	fmt.Println(BigInt("20^3"))  // 8000
	fmt.Println(BigInt("2^3^3")) // 134217728
	fmt.Println(BigInt("2^27"))  // 134217728
	nums := []string{
		"1",
		"123456",
		"123456789",
		"0x1",
		"0X1",
		"0o1",
		"0O1",
		"0b1",
		"0B1",
		"0xfffff",
		"aaa",
	}
	for _, s := range nums {
		log.Println(BigInt([]byte(s)))
	}

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
}
