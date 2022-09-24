package main

import (
	"BTC_Simulate/blockchain"
	"fmt"
)

func main() {
	blockChain := blockchain.NewBlockChain()
	blockChain.AddBlock("我向xx转了50个BTC")
	blockChain.AddBlock("xx向我转了50个BTC")

	for i, block := range blockChain.Blocks {
		fmt.Printf("===============当前区块高度：%d===============\n", i)
		fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：%x\n", block.CurHash)
		fmt.Printf("区块数据：%s\n", block.Data)
	}
}
