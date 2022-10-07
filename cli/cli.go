package cli

import (
	"BTC_Simulate/blockchain"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	blockchain *blockchain.Blockchain
}

const Usage = `
	newWallet                             "new a wallet"
	listAddresses                         "list all addresses"
	printBlockchain                       "print all blockchain data"
	getBalance --address USER             "get user balance"
	transfer FROM TO AMOUNT MINER DATA    "transfer"
`

func NewCLI(blockchain *blockchain.Blockchain) *CLI {
	return &CLI{blockchain}
}

func (cli *CLI) Run() {
	// 得到所有的命令
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}

	// 分析命令
	arg := args[1]
	switch arg {
	case "newWallet":
		fmt.Printf("新建钱包...\n")
		if len(args) == 2 {
			cli.NewWallet()
		} else {
			fmt.Printf("newWallet参数错误，请重新输入！\n")
			fmt.Printf(Usage)
		}
		break

	case "listAddresses":
		fmt.Printf("列出钱包地址...\n")
		if len(args) == 2 {
			cli.ListAddresses()
		} else {
			fmt.Printf("listAddresses参数错误，请重新输入！\n")
			fmt.Printf(Usage)
		}
		break

	case "printBlockchain":
		fmt.Printf("打印区块链...\n")
		if len(args) == 2 {
			cli.PrintBlockchain()
		} else {
			fmt.Printf("printBlockchain参数错误，请重新输入！\n")
			fmt.Printf(Usage)
		}
		break

	case "getBalance":
		fmt.Printf("读取余额...\n")
		if len(args) == 4 && args[2] == "--address" {
			user := args[3]
			cli.getBalance(user)
		} else {
			fmt.Printf("addBlock参数错误，请重新输入！\n")
			fmt.Printf(Usage)
		}
		break

	case "transfer":
		fmt.Printf("转账开始...\n")
		if len(args) == 7 {
			form := args[2]
			to := args[3]
			amount, _ := strconv.ParseFloat(args[4], 64)
			miner := args[5]
			data := args[6]
			cli.transfer(form, to, amount, miner, data)
		} else {
			fmt.Printf("transfer参数错误，请重新输入！\n")
			fmt.Printf(Usage)
		}
		break

	default:
		fmt.Printf("命令错误，请重新输入！\n")
		fmt.Printf(Usage)
		break
	}
}
