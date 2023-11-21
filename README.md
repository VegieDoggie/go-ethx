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
> 默认情况下，并发的最大数量=rpc的数量，这一限制是安全的，你可以放心地循环或并发地进行接口调用。
```go
go get github.com/VegieDoggie/go-ethx@v1.5.6
```
> 1- 可靠接口请求(完整接口请查看文档末尾的`接口概览`)
```go
func main() {
    // 说明: Clientx初始化时会自动剔除其中不连通的RPC接口
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
    // 可靠地合约接口请求(适用于读取链上合约数据)
    must := clientx.NewMust(UniswapV2Pair.NewUniswapV2Pair, "0xbe8561968ce5f9a9bf5cf6a117dfdee1b0e56d75")
    // 方式一
    token0 := must("Token0").(common.Address)
    token1 := must("Token1").(common.Address)
    reserves := ethx.Type[R](must("GetReserves"))
    log.Println(token0, token1, reserves)
    
    // 方式二
    token0 = must(UniswapV2Pair.UniswapV2Pair.Token0).(common.Address)
    token1 = must(UniswapV2Pair.UniswapV2Pair.Token1).(common.Address)
    reserves = ethx.Type[R](must(UniswapV2Pair.UniswapV2Pair.GetReserves))
    log.Println(token0, token1, reserves)
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
    
为了确保同一个区块范围同时被多个rpc节点扫描，这样即使某节点出问题，也不会漏块。

**假如你正在持续扫描最新块的日志**，在`延迟=3，回溯=800`的情况下，同一个区块预计将被扫描 `800/3=266次`，这是恐怖的数量!

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
fmt.Println(BigInt("1e3"))   // 1000
fmt.Println(BigInt("10E3"))  // 10000
fmt.Println(BigInt("2^3"))   // 8
fmt.Println(BigInt("20^3"))  // 8000
fmt.Println(BigInt("2^3^3")) // 134217728
fmt.Println(BigInt("2^27"))  // 134217728
// 随机
fmt.Println(Hash(RandBytes(32)))    // 0x06e4ba3da81342545e60108a576ef5590ee56800ef285bd692923b696f05fa44
fmt.Println(Address(RandBytes(20))) // 0x8807B00e663fD283ff7e9C1291EFF9D6963290Da

// Rpc 检查连通性和响应速率
raw := "https://bsc-dataseed1.defibit.io\t29599361\t1.263s\t\t\nConnect Wallet\nhttps://bsc-dataseed4.binance.org\t29599361\t1.263s\t\t\nConnect Wallet\nhttps://bsc.blockpi.network/v1/rpc/public\t29599361\t1.263s\t\t\nConnect Wallet\nhttps://bsc-dataseed2.defibit.io\t29599361\t1.264s\t\t\nConnect Wallet\nhttps://bsc-dataseed1.binance.org\t29599361\t1.264s\t\t\nConnect Wallet\nhttps://koge-rpc-bsc.48.club\t29599361\t1.264s\t\t\nConnect Wallet\nhttps://bsc-dataseed3.binance.org\t29599361\t1.265s\t\t\nConnect Wallet\nhttps://bsc-dataseed4.defibit.io\t29599361\t1.548s\t\t\nConnect Wallet\nhttps://bsc-dataseed2.binance.org\t29599361\t1.548s\t\t\nConnect Wallet\nhttps://bsc-mainnet.rpcfast.com?api_key=xbhWBI1Wkguk8SNMu1bvvLurPGLXmgwYeC4S6g2H7WdwFigZSmPWVZRxrskEQwIf\t29599361\t1.548s\t\t\nConnect Wallet\nhttps://bsc-dataseed3.ninicoin.io\t29599361\t1.549s\t\t\nConnect Wallet\nhttps://bsc-mainnet.gateway.pokt.network/v1/lb/6136201a7bad1500343e248d\t29599361\t1.549s\t\t\nConnect Wallet\nhttps://bsc-dataseed2.ninicoin.io\t29599361\t1.840s\t\t\nConnect Wallet\nhttps://rpc-bsc.48.club\t29599361\t1.843s\t\t\nConnect Wallet\nhttps://rpc.ankr.com/bsc\t29599361\t1.843s\t\t\nConnect Wallet\nhttps://1rpc.io/bnb\t29599361\t1.843s\t\t\nConnect Wallet\nhttps://bsc-dataseed3.defibit.io\t29599361\t1.844s\t\t\nConnect Wallet\nhttps://bsc.publicnode.com\t29599361\t1.844s\t\t\nConnect Wallet\nhttps://bsc-dataseed4.ninicoin.io\t29599361\t2.115s\t\t\nConnect Wallet\nhttps://endpoints.omniatech.io/v1/bsc/mainnet/public\t29599361\t2.115s\t\t\nConnect Wallet\nhttps://bsc-mainnet.nodereal.io/v1/64a9df0874fb4a93b9d0a3849de012d3\t29599361\t2.115s\t\t\nConnect Wallet\nhttps://binance.nodereal.io\t29599361\t2.387s\t\t\nConnect Wallet\nhttps://bsc-mainnet.public.blastapi.io\t29599361\t2.387s\t\t\nConnect Wallet\nhttps://bsc.meowrpc.com\t29599361\t2.654s\t\t\nConnect Wallet\nhttps://bscrpc.com\t29599360\t1.549s\t\t\nConnect Wallet\nhttps://bnb.api.onfinality.io/public\t29599359\t1.845s\t\t\nConnect Wallet\nhttps://bsc-dataseed1.ninicoin.io\t29598929\t1.814s\t\t\nConnect Wallet\nwss://bsc-ws-node.nariox.org\t\t\t\t\nhttps://nodes.vefinetwork.org/smartchain\t\t\t\t\nConnect Wallet\nhttps://bsc.rpcgator.com"
reliableRpcList, _, _, _ := CheckRpcConn(raw)
CheckRpcSpeed(reliableRpcList...)
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
# Rpc 检查连通性和响应速率
CheckRpcSpeed(rpcLike ...string) (rpcSpeedMap map[string]time.Duration)
CheckRpcConn(rpcLike ...string) (reliableRpcList []string, badRpcList []string, reliableClients []*ethclient.Client, reliableRpcMap map[*ethclient.Client]string)

```
# 开发计划

对比SubGraph: 
- 子图的同步时间过长，动辄便需要数个小时
- 子图有一定的开发门槛，需要学习搭建和数据定义
- Go-ETHx对原始日志的同步时间往往只需要几分钟(扫块间隔越大，RPC数量越多，时间越短! 理论上没有上限，可无限升级)
- Go-ETHx的原始日志解析完全采用Go语言，没有门槛

初步计划(近期时间和精力有限，最终实现时间不定): 
- 计划集成ERC20转账，归集，分发，查税点/貔貅等功能，时间不定
- 计划剥离所有的go-etherum依赖，时间不定，会大重构，发布新仓库
- 计划兼容所有EVM链及其扩展链，比如: CFX，时间不定，最终形态的目标是全链兼容，实现0成本/低成本迁移
- 计划建立Solidity常见内置函数的工具集(除了ts版的viem外，很少有名称对应的库支持，使用起来很不方便)
- 计划增加签名校验，尤其适合中心化交易所的部分
- 考虑提供类似子图的功能，短期内暂时没有时间
- 考虑开发ts, java，rust等主流语言版本，短期内暂时没有时间
