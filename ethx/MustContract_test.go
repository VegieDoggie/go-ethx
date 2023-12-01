package ethx

import (
	"fmt"
	"github.com/VegieDoggie/go-ethx/internal/TestLog"
	"math/big"
	"testing"
)

var consts = struct {
	TestLog *TestLog.TestLog
}{
	TestLog: new(TestLog.TestLog),
}

func Test_Call0_localhost(t *testing.T) {
	//// 0x5FbDB2315678afecb367f032d93F642f64180aa3
	//clientx := NewSimpleClientx([]string{"http://127.0.0.1:8545/"}, 1)
	//
	//testLog, _ := TestLog.NewTestLog(Address("0x5FbDB2315678afecb367f032d93F642f64180aa3"), clientx.NextClient())
	//log.Println(testLog.GetData0(nil))
}

func Test_Call0(t *testing.T) {
	// polygonMumbai 0xF51C6A832bF823C571456927A1751D98e5593230
	clientx := NewSimpleClientx([]string{"https://rpc-mumbai.maticvigil.com"}, 2)
	testLogAddr := Address("0xF51C6A832bF823C571456927A1751D98e5593230")
	testLog, _ := TestLog.NewTestLog(testLogAddr, clientx.NextClient())
	data0, _ := testLog.GetData0(nil)
	if data0.Cmp(BigInt(100)) != 0 {
		panic(fmt.Errorf("GetData0"))
	}

	config := NewClientxConfig()
	config.Event.IntervalBlocks = 20
	config.Event.OverrideBlocks = 20
	config.Event.DelayBlocks = 2
	mustTestLog := clientx.NewMustContract(TestLog.NewTestLog, testLogAddr, config)
	if mustTestLog.Read0(consts.TestLog.GetData0).(*big.Int).Cmp(BigInt(100)) != 0 {
		panic(fmt.Errorf("GetData0"))
	}
	//log.Println(mustTestLog.Read(consts.TestLog.GetData1, 1))
	//log.Println(mustTestLog.Read(consts.TestLog.GetData2, 1, 2))
}

func Test_Subscribe(t *testing.T) {
	//// 0x5FbDB2315678afecb367f032d93F642f64180aa3
	//clientx := NewSimpleClientx([]string{"http://127.0.0.1:8545/"}, 10)
	//config := EventConfig{
	//	IntervalBlocks: 20,
	//	OverrideBlocks: 10,
	//	DelayBlocks:    2,
	//}
	//mustTestLog := clientx.NewMustContract(TestLog.NewTestLog, "0x5FbDB2315678afecb367f032d93F642f64180aa3", config)
	//logCh := make(chan *TestLog.TestLogIndex0, 128)
	//sub := mustTestLog.Subscribe(logCh, 0)
	//sub.Err()
	//for i := 0; i < 105; i++ {
	//	l := <-logCh
	//	log.Println(l.Raw.TxHash, l.Raw.BlockNumber)
	//	//if i == 30 {
	//	//	sub.Unsubscribe()
	//	//	return
	//	//}
	//}
}

func Test_Subscribe_index(t *testing.T) {
	//// 0x5FbDB2315678afecb367f032d93F642f64180aa3
	//clientx := NewSimpleClientx([]string{"http://127.0.0.1:8545/"}, 10)
	//config := EventConfig{
	//	IntervalBlocks: 20,
	//	OverrideBlocks: 10,
	//	DelayBlocks:    0,
	//}
	//mustTestLog := clientx.NewMustContract(TestLog.NewTestLog, "0x5FbDB2315678afecb367f032d93F642f64180aa3", config)
	//logCh := make(chan *TestLog.TestLogIndex1, 128)
	////sub := mustTestLog.Subscribe(logCh, 0)
	//sub := mustTestLog.Subscribe(logCh, 0, nil)
	//sub.Err()
	//for i := 0; i < 105; i++ {
	//	l := <-logCh
	//	log.Println(l.Raw.TxHash, l.Raw.BlockNumber)
	//	//if i == 30 {
	//	//	sub.Unsubscribe()
	//	//	return
	//	//}
	//}
}

func Test_NewTestLog_FilterIndex0(t *testing.T) {
	//clientx := NewSimpleClientx([]string{"http://127.0.0.1:8545/"}, 10)
	//testLog, _ := TestLog.NewTestLog(Address("0x5FbDB2315678afecb367f032d93F642f64180aa3"), clientx.NextClient())
	//end := uint64(2000)
	//index0Iterator, err := testLog.FilterIndex1(&bind.FilterOpts{
	//	Start:   0,
	//	End:     &end,
	//	Context: nil,
	//}, []common.Address{Address("0x0000000000000000000000000000000000000001")})
	//if err != nil {
	//	panic(err)
	//}
	//for i := 0; i < 100; i++ {
	//	if index0Iterator.Next() {
	//		log.Println(index0Iterator.Event.Raw.TxHash, index0Iterator.Event.Raw.BlockNumber)
	//	}
	//}
}
