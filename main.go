package main

import (
	"BTC_Simulate/blockchain"
	"BTC_Simulate/cli"
)

func main() {
	newBlockchain := blockchain.NewBlockchain("Satoshi Nakamoto")
	newCli := cli.NewCLI(newBlockchain)
	newCli.Run()
}
