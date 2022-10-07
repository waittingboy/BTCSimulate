package blockchain

import (
	"BTC_Simulate/wallet"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"log"
)

type Transaction struct {
	TXId      []byte     // 交易ID
	TXInputs  []TXInput  // 交易输入数组
	TXOutputs []TXOutput // 交易输出数组
}

type TXInput struct {
	TXId      []byte // 引用的交易ID
	Index     int    // 引用output在交易中的索引值
	Signature []byte // 数字签名
	PublicKey []byte // 公钥
}

type TXOutput struct {
	Amount        float64 // 转账金额
	PublicKeyHash []byte  // 公钥Hash
}

const reward = 6.25

// 创建TXOutput
func NewTXOutput(amount float64, address string) *TXOutput {
	output := TXOutput{
		Amount: amount,
	}

	output.PublicKeyHash = GetPublicKeyHash(address)

	return &output
}

// 计算公钥Hash
func GetPublicKeyHash(address string) []byte {
	// 解码
	addressByte := base58.Decode(address)

	// 计算公钥Hash
	publicKeyHash := addressByte[1 : len(addressByte)-4]

	return publicKeyHash
}

// 设置交易ID
func (transaction *Transaction) SetTXId() {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(transaction)
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
	// 构建Input：矿工挖矿时无需指定公钥，所以PublicKey字段可以随意填写
	input := TXInput{[]byte{}, -1, nil, []byte(data)}
	// 构建Output
	output := NewTXOutput(reward, miner)

	transaction := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
	transaction.SetTXId()

	return &transaction
}

// 创建普通交易
func NewTransaction(from, to string, amount float64, blockchain *Blockchain) *Transaction {
	wallets := wallet.NewWallets()
	fromWallet := wallets.WalletsMap[from]
	if fromWallet == nil {
		fmt.Printf("from地址不存在！\n")
		return nil
	}

	//privateKey := fromWallet.PrivateKey
	publicKey := fromWallet.PublicKey

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
			input := TXInput{[]byte(id), index, nil, append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)}
			inputs = append(inputs, input)
		}
	}

	// 构建Output
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)

	if totalAmount > amount {
		output = NewTXOutput(totalAmount-amount, from)
		outputs = append(outputs, *output)
	}

	transaction := Transaction{[]byte{}, inputs, outputs}
	transaction.SetTXId()

	return &transaction
}
