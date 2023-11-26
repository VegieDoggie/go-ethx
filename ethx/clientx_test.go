package ethx

import (
	"log"
	"testing"
)

func Test_Clientx(t *testing.T) {
	reliableList := []string{
		"https://data-seed-prebsc-2-s2.bnbchain.org:8545",
		"https://endpoints.omniatech.io/v1/bsc/testnet/public",
		"https://bsc-testnet.publicnode.com",
		//"https://data-seed-prebsc-2-s3.bnbchain.org:8545",
		//"https://bsc-testnet.public.blastapi.io",
		//"https://data-seed-prebsc-1-s3.bnbchain.org:8545",
		//"https://data-seed-prebsc-1-s1.bnbchain.org:8545",
		//"https://data-seed-prebsc-1-s2.bnbchain.org:8545",
		//"https://bsc-testnet.blockpi.network/v1/rpc/public",
		//"https://data-seed-prebsc-2-s1.bnbchain.org:8545",
		//"wss://bsc-testnet.publicnode.com",
	}
	clientx := NewSimpleClientx(reliableList)
	prikey := ""
	to := "0x87393E5971A58952D94C7b200663fB49138fF001"
	tx, receipt, err := clientx.Transfer(prikey, to, BigInt("1e18"))
	log.Println(err)
	log.Println(tx)
	log.Println(receipt)
}

func Test_Error(t *testing.T) {
	//queryTicker := time.NewTicker(time.Second)
	//defer queryTicker.Stop()
	//go func() {
	//	panic("Mock error")
	//}()
	//for i := 0; i < 10; i++ {
	//	log.Println(i)
	//	<-queryTicker.C
	//}
}
