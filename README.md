# Go-ETHx

面向 Dapp 应用开发的 Ethereum 扩展接口。

## 预览

- 保证可靠的接口请求功能
- 任意区间的区块日志扫描功能
- 常用的工具函数

## 快速开始

> 默认情况下，并发的最大数量=rpc的数量，这一限制是安全的，你可以放心地循环或并发地进行接口调用。

```go
go get github.com/VegieDoggie/go -ethx@v1.6.1

rpcList := []string{
"https://data-seed-prebsc-2-s2.binance.org:8545",
"https://data-seed-prebsc-2-s1.binance.org:8545",
...
}
// 可靠的接口请求: 客户端
clientx := ethx.NewSimpleClientx(rpcList)
tx, receipt, err := clientx.Transfer(prikey, to, BigInt("1e18"))

// 可靠的接口请求: 合约
blockNumber := clientx.NewMust(UniswapV2Pair.NewUniswapV2Pair, "0xbe8561968ce5f9a9bf5cf6a117dfdee1b0e56d75")
token0 := must("Token0").(common.Address) // = must(UniswapV2Pair.UniswapV2Pair.Token0).(common.Address)
token1 := must("Token1").(common.Address) // = must(UniswapV2Pair.UniswapV2Pair.Token1).(common.Address)

// 任意区间的区块日志扫描功能
var topics [][]common.Hash
var addresses []common.Address
scanner := clientx.NewScanner(topics, addresses, 200, 800, 3)
logs, _ := scanner.Scan(0, clientx.BlockNumber())

// 常用工具
b0 := ethx.BigInt(0)
b1 := ethx.BigInt("0xbe8561968ce5f9a9bf5cf6a117dfdee1b0e56d75")
b2 := ethx.BigInt("1e18") // *big.Int: 1_000_000_000_000_000_000

addr0 := ethx.Address(0)
addr1 := ethx.Address("0xbe8561968ce5f9a9bf5cf6a117dfdee1b0e56d75")
addr2 := ethx.Address(prikey, true)

sum0 := ethx.Add(33, 100, "21") // *big.Int: 33+100+21
sum1 := ethx.Add([]int{33, 100}, 21) // *big.Int: 33+100+21

mul0 := ethx.Mul(33, "100", 2) // *big.Int: 33*100*21
mul1 := ethx.Mul([]int{33, "100", 2}) // *big.Int: 33*100*21
```
# 未来可能

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
