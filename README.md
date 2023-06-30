# Go-ETHx

面向 Dapp 应用开发的 Ethereum 扩展接口。

## 背景

在开发Dapp应用时，部分功能只能在中心端实现，于是最大的问题出现了: 如何保证中心化系统及时地探测到去中心端发出的事件并且保证不丢块? 为了解决这个问题，区块日志扫描功能便诞生了! 为了保证扫块机制的可靠性，除了引入节点轮换的机制外，还特意增加了延迟扫块和重复扫块的机制，从而保证扫块的又快又稳。

除了扫块以外，中心端↔去中心端的转账也需要解决可靠性问题，核心在于交易块的出块监听，这个问题也得到了很好的解决。

此外，还有Dapp开发中的常见场景，比如数据库和传输得到的地址，哈希字符串的转换，也做了一些优化。


## 预览

- 提供保证可靠的接口请求功能，自动轮换请求接口，直到请求成功
- 任意区间的区块日志扫描功能，自动划分区间拉取完整的日志
- 常用的转换函数及大数运算函数

## 快速开始

```go
go get github.com/VegetableDoggies/go-ethx@v2.0.0
```
> 1- 可靠接口请求(完整接口请查看文档末尾的`接口概览`)
```go
func main() {
    rpcList := []string{
        "https://data-seed-prebsc-2-s2.binance.org:8545",
        "https://data-seed-prebsc-2-s1.binance.org:8545",
    }
    weights := []int{1, 2}
    // 参数3: 交易默认确认块数
    clientx := NewClientx(rpcList, weights,30)

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
// 控制台输出
AutoBlockNumber 31138256
AutoBlockNumber 31138256
AutoBlockNumber 31138257
AutoBlockNumber 31138257
AutoBlockNumber 31138257
AutoBlockNumber 31138258
AutoBlockNumber 31138258
AutoBlockNumber 31138258
AutoBlockNumber 31138259
AutoBlockNumber 31138259
......
```
> 2- 区块日志扫描功能
------
解答1: 为什么要延迟扫块? 
    
为了防止扫描到叔块上的日志，这属于假交易，很危险

------

解答2: 为什么要回溯重复扫块? 
    
为了确保同一个区块范围同时被多个rpc节点扫描，这样即使某节点出问题，也不会漏块。假如你正在持续扫描最新块的日志，在`延迟=3，回溯=800`的情况下，同一个区块预计将被扫描 `800/3=266次`，这是恐怖的数量!

------
```go
func main() {
    rpcList := []string{
        "https://data-seed-prebsc-2-s2.binance.org:8545",
        "https://data-seed-prebsc-2-s1.binance.org:8545",
    }
    clientx := NewClientx(rpcList, nil, 30)
    // intervalBlocks=200: 最终被切分后的扫块范围
    // overrideBlocks=800: 每次回溯重复扫的范围，若初始 from=800，则真实的 from=0
    // delayBlocks=3:   延迟扫块的数量，若初始 to=1000，则真实的 to=997
    scanner := clientx.NewScanner(nil, nil, 200, 800, 3)
    logs, _ := scanner.Scan(0, 3000)
    fmt.Println(logs)
}
// 控制台输出
2023/06/30 15:34:10 Scan start for: 1200-1399 Success!
2023/06/30 15:34:10 Scan start for: 2800-2997 Success!
2023/06/30 15:34:10 Scan start for: 0-199 Success!
2023/06/30 15:34:11 Scan start for: 200-399 Success!
2023/06/30 15:34:12 Scan start for: 1800-1999 Success!
2023/06/30 15:34:12 Scan start for: 400-599 Success!
2023/06/30 15:34:12 Scan start for: 600-799 Success!
2023/06/30 15:34:13 Scan start for: 800-999 Success!
2023/06/30 15:34:13 Scan start for: 1000-1199 Success!
2023/06/30 15:34:14 Scan start for: 1400-1599 Success!
2023/06/30 15:34:14 Scan start for: 1600-1799 Success!
2023/06/30 15:34:15 Scan start for: 2200-2399 Success!
2023/06/30 15:34:16 Scan start for: 2000-2199 Success!
2023/06/30 15:34:16 Scan start for: 2400-2599 Success!
2023/06/30 15:34:17 Scan start for: 2600-2799 Success!
[{0x1B0362f94487f2e8a6cDC2b9D352bd63F147eBc5 [0x8be0079c5316591......
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
## 接口概览
> package: `ethx`
```cmd
# -----------------------------------
#              可靠客户端
# -----------------------------------
Clientx
NewClientx(rpcList []string, weights []int, notFoundBlocks uint64, limiter ...*rate.Limiter) *Clientx
    BlockNumber() (blockNumber uint64)
    ChainID() (chainID *big.Int)
    NetworkID() (networkID *big.Int)
    BalanceAt(account any, blockNumber any) (balance *big.Int)
    PendingBalanceAt(account any) (balance *big.Int)
    NonceAt(account any, blockNumber any) (nonce uint64)
    PendingNonceAt(account any) (nonce uint64)
    FilterLogs(q ethereum.FilterQuery) (logs []types.Log)
    SuggestGasPrice() (gasPrice *big.Int)
    SuggestGasTipCap() (gasTipCap *big.Int)
    FeeHistory(blockCount uint64, lastBlock any, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory)
    StorageAt(account any, keyHash any, blockNumber any) (storage []byte)
    PendingStorageAt(account any, keyHash any) (storage []byte)
    CodeAt(account any, blockNumber any) (code []byte)
    PendingCodeAt(account any) (code []byte)
    BlockByHash(hash any, notFoundBlocks ...uint64) (block *types.Block, err error)
    BlockByNumber(blockNumber any, notFoundBlocks ...uint64) (block *types.Block, err error)
    HeaderByHash(hash any, notFoundBlocks ...uint64) (header *types.Header, err error)
    HeaderByNumber(blockNumber any, notFoundBlocks ...uint64) (header *types.Header, err error)
    TransactionByHash(hash any, notFoundBlocks ...uint64) (tx *types.Transaction, isPending bool, err error)
    TransactionSender(tx *types.Transaction, blockHash any, index uint, notFoundBlocks ...uint64) (sender common.Address, err error)
    TransactionCount(blockHash any, notFoundBlocks ...uint64) (count uint, err error)
    PendingTransactionCount() (count uint)
    TransactionInBlock(blockHash any, index uint, notFoundBlocks ...uint64) (tx *types.Transaction, err error)
    TransactionReceipt(txHash any, notFoundBlocks ...uint64) (receipt *types.Receipt, err error)
    # 等待交易确认
    WaitMined(tx *types.Transaction, confirmBlocks uint64, notFoundBlocks ...uint64) (*types.Receipt, error)
    # 等待部署确认
    WaitDeployed(tx *types.Transaction, confirmBlocks uint64, notFoundBlocks ...uint64) (common.Address, error)
    # 新建日志扫描器: Scanner
    NewScanner(topics [][]common.Hash, addresses []common.Address, intervalBlocks uint64, overrideBlocks uint64, delayBlocks uint64) *Scanner
    Scanner
        Scan(from uint64, to uint64) (logs []types.Log, addressTopicLogsMap AddressTopicLogsMap)
# -----------------------------------
#              常用工具
# -----------------------------------
# 任意数转common.Hash
Hash(hashLike any) common.Hash
# 任意数转common.Address
Address(addressLike any) common.Address
# 任意数转*big.Int
BigInt(numLike any) *big.Int
# 任意数转uint64
Uint64(numLike any) uint64
# 任意数转Int64
Int64(numLike any) int64
# 任意数的基本运算: 加减乘除 大小比较
Add(numLike ...any) *big.Int
Sub(numBase any, numLike ...any) *big.Int
Mul(numLike ...any) *big.Int
Div(numLike0 any, numLike1 any) *big.Int
Sum(numLike ...any) *big.Int
Gte(numLike0 any, numLike1 any, isAbs ...bool) bool
Gt(numLike0 any, numLike1 any, isAbs ...bool) bool
Lte(numLike0 any, numLike1 any, isAbs ...bool) bool
Lt(numLike0 any, numLike1 any, isAbs ...bool) bool
# 基础判断
Is0x(s string) bool
Is0b(s string) bool
Is0o(s string) bool
# 随机字节组，如例: ethx.Hash(ethx.RandBytes(32))
RandBytes(len int) []byte
```
