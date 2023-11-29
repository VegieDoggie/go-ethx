package ethx

import (
	"github.com/VegieDoggie/go-ethx/internal/TestLog"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"log"
	"testing"
	"time"
)

func TestTest(t *testing.T) {
	log.Println(segmentCallback(0, 100, EventConfig{
		IntervalBlocks: 10,
		OverrideBlocks: 0,
		DelayBlocks:    0,
	}, func(from, to uint64) {
		log.Println(from, to)
	}))
	select {}
}

func Test_NewMustContract_Subscribe(t *testing.T) {
	// 0x5FbDB2315678afecb367f032d93F642f64180aa3
	clientx := NewSimpleClientx([]string{"http://127.0.0.1:8545/"}, 10)
	clientx.miningInterval = 3 * time.Second
	config := EventConfig{
		IntervalBlocks: 20,
		OverrideBlocks: 0,
		DelayBlocks:    0,
	}
	mustTestLog := clientx.NewMustContract(TestLog.NewTestLog, "0x5FbDB2315678afecb367f032d93F642f64180aa3", config)
	logCh := make(chan *TestLog.TestLogIndex0, 128)
	sub := mustTestLog.Subscribe(0, logCh)
	for i := 0; i < 100; i++ {
		l := <-logCh
		log.Println(l.Raw.TxHash, l.Raw.BlockNumber)
		if i == 30 {
			sub.Unsubscribe()
			return
		}
	}
}

func Test_NewTestLog_FilterIndex0(t *testing.T) {
	clientx := NewSimpleClientx([]string{"http://127.0.0.1:8545/"}, 10)
	testLog, _ := TestLog.NewTestLog(Address("0x5FbDB2315678afecb367f032d93F642f64180aa3"), clientx.NextClient())
	end := uint64(2000)
	index0Iterator, err := testLog.FilterIndex0(&bind.FilterOpts{
		Start:   0,
		End:     &end,
		Context: nil,
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		if index0Iterator.Next() {
			log.Println(index0Iterator.Event.Raw.TxHash, index0Iterator.Event.Raw.BlockNumber)
		}
	}
}
