package cli

import (
	"BTC_Simulate/blockchain"
	"BTC_Simulate/utils"
	"BTC_Simulate/wallet"
	"fmt"
)

func (cli *CLI) NewWallet() {
	wallets := wallet.NewWallets()
	address := wallets.AddWallet()
	fmt.Printf("地址：%s\n", address)
}

func (cli *CLI) ListAddresses() {
	wallets := wallet.NewWallets()
	addresses := wallets.ListAllAddresses()
	for _, address := range addresses {
		fmt.Printf("地址：%s\n", address)
	}
}

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
			fmt.Printf("当前区块第%d个交易的ID：%x\n", i, transaction.TXId)
			fmt.Println(transaction.ToString())
		}

		if len(block.BlockHead.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(user string) {
	if !wallet.IsValidAddress(user) {
		fmt.Printf("地址无效，请重新输入！\n")
		return
	}

	UTXOs := cli.blockchain.FindUTXOs(user)
	balance := 0.0

	for _, UTXO := range UTXOs {
		balance += UTXO.Amount
	}

	fmt.Printf("地址%s的余额为：%f\n", user, balance)
}

func (cli *CLI) transfer(from, to string, amount float64, miner, data string) {
	if !wallet.IsValidAddress(from) {
		fmt.Printf("from地址无效，请重新输入！\n")
		return
	}

	if !wallet.IsValidAddress(to) {
		fmt.Printf("to地址无效，请重新输入！\n")
		return
	}

	if !wallet.IsValidAddress(miner) {
		fmt.Printf("miner地址无效，请重新输入！\n")
		return
	}

	// 创建挖矿交易
	coinbase := blockchain.NewCoinbase(miner, data)
	// 创建普通交易
	transaction := blockchain.NewTransaction(from, to, amount, cli.blockchain)
	if transaction == nil {
		return
	}
	// 将交易添加到区块中
	cli.blockchain.AddBlock([]*blockchain.Transaction{coinbase, transaction})
}
