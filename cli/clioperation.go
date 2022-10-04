package cli

import (
	"BTC_Simulate/blockchain"
	"BTC_Simulate/utils"
	"fmt"
)

func (cli *CLI) PrintBlockchain() {
	iterator := cli.blockchain.NewIterator()
	for {
		block := iterator.Next()
		fmt.Printf("========================================\n")
		fmt.Printf("当前版本号：%d\n", block.BlockHead.Version)
		fmt.Printf("前区块哈希值：%x\n", block.BlockHead.PrevHash)
		fmt.Printf("梅克尔根：%x\n", block.BlockHead.MerkelRoot)
		fmt.Printf("区块时间：%s\n", utils.TimeFormat(block.BlockHead.TimeStamp))
		fmt.Printf("随机数：%d\n", block.BlockHead.Nonce)
		fmt.Printf("当前区块哈希值：%x\n", block.CurHash)
		for i, transaction := range block.Transactions {
			fmt.Printf("当前区块第%d个交易的ID为：%x\n", i, transaction.TXId)
		}
		fmt.Printf("区块数据：%s\n", block.Transactions[0].TXInputs[0].Sig)

		if len(block.BlockHead.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(user string) {
	UTXOs := cli.blockchain.FindUTXOs(user)
	balance := 0.0

	for _, UTXO := range UTXOs {
		balance += UTXO.Amount
	}

	fmt.Printf("%s的余额为：%f\n", user, balance)
}

func (cli *CLI) transfer(from, to string, amount float64, miner, data string) {
	// 创建挖矿交易
	coinbase := blockchain.NewCoinbase(miner, data)
	// 创建普通交易
	transaction := blockchain.NewTransactionForSingle(from, to, amount, cli.blockchain)
	if transaction == nil {
		fmt.Printf("转账失败！\n")
		return
	}
	// 将交易添加到区块中
	cli.blockchain.AddBlock([]*blockchain.Transaction{coinbase, transaction})
}
