# Go-ETHx

面向 Dapp 应用开发的 Ethereum 扩展接口。

## 预览

- 提供保证可靠的接口请求功能，自动轮换请求接口，直到请求成功
- 任意区间的区块日志扫描功能，自动划分区间拉取完整的日志
- 常用的转换函数及大数运算函数

## 快速开始

```cmd
go get github.com/VegetableDoggie/go-ethx@v1.0.0
```
> 1- 可靠接口请求
```go
// 默认权重值: 1
func main() {
	clientIterator := ethx.NewClientIterator([]string{
		"https://data-seed-prebsc-2-s2.binance.org:8545",
		"https://data-seed-prebsc-2-s1.binance.org:8545",
	})
	clientx := ethx.NewClientx(clientIterator, 15, 10)
	fmt.Println(clientx.BlockNumber(), clientx.LocalBlockNumber())
}

// 自定义权重值
func main() {
    clientIterator := ethx.NewClientIteratorWithWeight([]string{
    "https://data-seed-prebsc-2-s2.binance.org:8545",
    "https://data-seed-prebsc-2-s1.binance.org:8545",
    }, []int{3, 5})
    clientx := ethx.NewClientx(clientIterator, 15, 10)
    fmt.Println(clientx.BlockNumber(), clientx.LocalBlockNumber())
}
```
> 2- 区块日志扫描功能
```go
func main() {
    clientIterator := ethx.NewClientIterator([]string{
    "https://data-seed-prebsc-2-s2.binance.org:8545",
    "https://data-seed-prebsc-2-s1.binance.org:8545",
    })
    scanner := ethx.NewScanner(clientIterator, nil, nil)
    logs, _ := scanner.Scan(context.TODO(), 0, 2000)
    fmt.Println(logs)
}
```

> 3- 常用的转换函数及大数运算函数
```go
// To Hash
fmt.Println(ethx.Hash("1")) // 0x0000000000000000000000000000000000000000000000000000000000000001
fmt.Println(ethx.Hash(1)) // 0x0000000000000000000000000000000000000000000000000000000000000001
fmt.Println(ethx.Hash("0x1")) // 0x0000000000000000000000000000000000000000000000000000000000000001
fmt.Println(ethx.Hash("0b1")) // 0x0000000000000000000000000000000000000000000000000000000000000001
fmt.Println(ethx.Hash("0o1")) // 0x0000000000000000000000000000000000000000000000000000000000000001

// To Address
fmt.Println(ethx.Address("1")) // 0x0000000000000000000000000000000000000001
fmt.Println(ethx.Address(1)) // 0x0000000000000000000000000000000000000001
fmt.Println(ethx.Address("0x1")) // 0x0000000000000000000000000000000000000001
fmt.Println(ethx.Address("0b1")) // 0x0000000000000000000000000000000000000001
fmt.Println(ethx.Address("0o1")) // 0x0000000000000000000000000000000000000001

// To BigInt
fmt.Println(ethx.BigInt("1")) // 1
fmt.Println(ethx.BigInt(1)) // 1
fmt.Println(ethx.BigInt("0x1")) // 1
fmt.Println(ethx.BigInt("0b1")) // 1
fmt.Println(ethx.BigInt("0o1")) // 1

// + - * / 和 比较
fmt.Println(ethx.Add("1", 2, 3, ethx.Hash(1)))    // 7
fmt.Println(ethx.Sum("1", 2, 3, ethx.Hash(1)))    // 7
fmt.Println(ethx.Sub(ethx.Address(1), ethx.Hash(0)))   // 1
fmt.Println(ethx.Mul("0x2", 5))  // 10
fmt.Println(ethx.Div("10", 2))   // 5
fmt.Println(ethx.Gte(ethx.Address(1), ethx.Address(2)))    // true
fmt.Println(ethx.Gt(ethx.Address(1), ethx.Address(1))) // false

// 随机
fmt.Println(ethx.Hash(ethx.RandBytes(32))) // 0x06e4ba3da81342545e60108a576ef5590ee56800ef285bd692923b696f05fa44
fmt.Println(ethx.Address(ethx.RandBytes(20))) // 0x8807B00e663fD283ff7e9C1291EFF9D6963290Da
```