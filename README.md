# Go-ETHx

面向 Dapp 应用开发的 Ethereum 扩展接口。

## 预览

- [快速开始](#快速开始): 如何创建**稳定可靠高并发**的 `*ethclient.Client`?
- [合约请求](#合约请求): 如何与**合约交互**? [可靠高并发]
- [日志订阅](#日志订阅): 如何进行`HTTP/WS`**订阅**? [可靠高并发]
- [日志扫描](#日志扫描): 如何**多合约多事件扫描**? [可靠高并发]
- [常用工具](#常用工具): 常用的**工具**函数，比如类型转换，加减乘除，私钥转换等

## 快速开始

> 默认情况下，每秒并发上限=rpc的数量(可自定义)，因此可以放心地循环或并发地进行接口调用。

```go
go get github.com/VegieDoggie/go-ethx

rpcList := []string{
"https://data-seed-prebsc-2-s2.binance.org:8545",
"https://data-seed-prebsc-2-s1.binance.org:8545",
...
}
// 可靠客户端(包含但不限于 *ethclient.Client 的方法)，如下: 
clientx := ethx.NewSimpleClientx(rpcList)

// 快捷转账/交易(支持携带数据)
tx, receipt, err := clientx.Transfer(prikey, to, BigInt("1e18"))
```
## 合约请求

```solidity
// 可靠的接口请求: 合约
uniswapV2Pair := clientx.NewMustContract(UniswapV2Pair.NewUniswapV2Pair, "0xbe8561968ce5f9a9bf5cf6a117dfdee1b0e56d75")
token0 := uniswapV2Pair.Read0("Token0").(common.Address) 
// same as: .Read0(new(UniswapV2Pair.UniswapV2Pair).Token0)
token1 := uniswapV2Pair.Read0("Token1").(common.Address)
// same as: .Read0(new(UniswapV2Pair.UniswapV2Pair).Token1)
```

## 日志订阅

> 支持任意起点，支持 HTTP/WSS

```solidity
mustMarketRouter := clientx.NewMustContract(MarketRouter.NewMarketRouter, m.contract)
ch := make(chan *MarketRouter.MarketRouterUpdateOrder)
// 解析chan自动监听: *MarketRouter.MarketRouterUpdateOrder
sub, newBlockNumberCh = mustMarketRouter.Subscribe(ch, startBlock) // newBlockNumberCh 返回值表示已经扫描到最新块
```

## 日志扫描

> 支持任意区间

```solidity
rpcs := []string{"https://rpc.arb1.arbitrum.gateway.fm", "https://rpc.ankr.com/arbitrum"...}
mabi, _ := Roles.RolesMetaData.GetAbi() // abigen Contract: <Roles>
topics := []common.Hash{
    mabi.Events["OwnershipTransferred"].ID,
    mabi.Events["RoleAdminChanged"].ID,
    mabi.Events["RoleGranted"].ID,
    mabi.Events["RoleRevoked"].ID,
}
clientx := ethx.NewSimpleClientx(rpcs, 4) // 每个RPC的最大并发数: 4 per/s
// 创建实例
config := ethx.NewClientxConfig()
config.Event.IntervalBlocks = 2000
config.Event.DelayBlocks = 0
config.Event.OverrideBlocks = 0
logger := clientx.NewRawLogger(nil, [][]common.Hash{topics}, config)
// 开启扫描
chLogs, _ := logger.Filter(0, 156099699)
for l := range chLogs {
    // TODO
}
```

## 常用工具

```solidity
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
