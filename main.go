package main

import (
	"BTC_Simulate/blockchain"
	"BTC_Simulate/cli"
)

func main() {
	newBlockchain := blockchain.NewBlockchain()
	newCli := cli.NewCLI(newBlockchain)
	newCli.Run()
}
