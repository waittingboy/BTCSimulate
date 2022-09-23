package main

import (
	"BTC_Simulate/blockchain"
	"fmt"
)

func main() {
	block := blockchain.NewBlock("我们的未来是星辰和大海！", []byte{})

	fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
	fmt.Printf("当前区块哈希值：%x\n", block.CurHash)
	fmt.Printf("区块数据：%s\n", block.Data)
}
