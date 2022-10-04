package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type Transaction struct {
	TXId      []byte     // 交易ID
	TXInputs  []TXInput  // 交易输入数组
	TXOutputs []TXOutput // 交易输出数组
}

type TXInput struct {
	TXId  []byte // 引用的交易ID
	Index int    // 引用output在交易中的索引值
	Sig   string // 解锁脚本：使用地址模拟
}

type TXOutput struct {
	Amount     float64 // 转账金额
	PubKeyHash string  // 锁定脚本：使用地址模拟
}

const reward = 6.25

// 设置交易ID
func (transaction *Transaction) SetTXId() {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(transaction)
	if err != nil {
		log.Panic("SetTXId：编码失败")
	}

	// 进行Hash运算
	hash := sha256.Sum256(buffer.Bytes())

	transaction.TXId = hash[:]
}

// 是否挖矿交易：只有一个Input，且Input引用的交易ID为空，引用的output的索引值为-1
func (transaction *Transaction) IsCoinbase() bool {
	inputs := transaction.TXInputs
	if len(inputs) == 1 && len(inputs[0].TXId) == 0 && inputs[0].Index == -1 {
		return true
	}

	return false
}

// 创建挖矿交易：只有一个Input和一个Output
func NewCoinbase(miner string, data string) *Transaction {
	// 构建Input：矿工挖矿时无需指定签名，所以sig字段可以随意填写
	input := TXInput{[]byte{}, -1, data}
	// 构建Output
	output := TXOutput{reward, miner}

	transaction := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	transaction.SetTXId()

	return &transaction
}

// 创建普通交易
func NewTransactionForSingle(from, to string, amount float64, blockchain *Blockchain) *Transaction {
	UTXOs, totalAmount := blockchain.FindNeedUTXOs(from, amount)
	if totalAmount < amount {
		fmt.Printf("余额不足！\n")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	// 构建Input
	for id, indexes := range UTXOs {
		for _, index := range indexes {
			input := TXInput{[]byte(id), index, from}
			inputs = append(inputs, input)
		}
	}

	// 构建Output
	output := TXOutput{amount, to}
	outputs = append(outputs, output)

	if totalAmount > amount {
		output := TXOutput{totalAmount - amount, from}
		outputs = append(outputs, output)
	}

	transaction := Transaction{[]byte{}, inputs, outputs}
	transaction.SetTXId()

	return &transaction
}
