package main

import (
	"BTC_Simulate/blockchain"
	"BTC_Simulate/cli"
)

func main() {
	newBlockchain := blockchain.NewBlockchain("1JbPyoZNo4gqtjrv2PuDdiEJMCvs2MVqfa")
	if newBlockchain == nil {
		return
	}
	newCli := cli.NewCLI(newBlockchain)
	newCli.Run()
}
