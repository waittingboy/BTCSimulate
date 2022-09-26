package cli

import (
	"BTC_Simulate/blockchain"
	"fmt"
	"os"
)

type CLI struct {
	blockchain *blockchain.Blockchain
}

const Usage = `
	addBlock --data DATA    "add block to blockchain"
	printBlockchain         "print all blockchain data"
`

func NewCLI(blockchain *blockchain.Blockchain) *CLI {
	return &CLI{blockchain}
}

func (cli *CLI) Run() {
	// 得到所有的命令
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("错误的命令或参数，请重新输入！")
		fmt.Printf(Usage)
		return
	}

	// 分析命令
	arg := args[1]
	switch arg {
	case "addBlock":
		if len(args) == 4 && args[2] == "--data" {
			data := args[3]
			cli.AddBlock(data)
		} else {
			fmt.Printf("addBlock参数错误，请重新输入！")
			fmt.Printf(Usage)
		}
		break

	case "printBlockchain":
		if len(args) == 2 {
			cli.PrintBlockchain()
		} else {
			fmt.Printf("printBlockchain参数错误，请重新输入！")
			fmt.Printf(Usage)
		}
 		break

	default:
		fmt.Printf("错误的命令或参数，请重新输入！")
		fmt.Printf(Usage)
		break
	}
}
