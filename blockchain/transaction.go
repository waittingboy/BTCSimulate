package blockchain

import (
	"BTC_Simulate/utils"
	"BTC_Simulate/wallet"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
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

	output.PublicKeyHash = utils.GetPublicKeyHash(address)

	return &output
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

	quotedTransactions := blockchain.FindQuotedTransactions(&transaction)
	transaction.Sign(wallets.WalletsMap[from].PrivateKey, quotedTransactions)

	return &transaction
}

// 交易签名
func (transaction *Transaction) Sign(privateKey *ecdsa.PrivateKey, quotedTransactions map[string]*Transaction) {
	if transaction.IsCoinbase() {
		return
	}

	transactionCopy := transaction.TrimmedCopy()

	for i, input := range transactionCopy.TXInputs {
		// 设置交易源数据Hash
		transactionCopy.SetDataHash(quotedTransactions, i, &input)

		// 进行数字签名
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, transactionCopy.TXId)
		if err != nil {
			log.Panic("交易签名失败！")
		}

		// 将数字签名赋值给transaction对应input的Signature
		transaction.TXInputs[i].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

// 交易验证
func (transaction *Transaction) Verify(quotedTransactions map[string]*Transaction) bool {
	transactionCopy := transaction.TrimmedCopy()

	for i, input := range transactionCopy.TXInputs {
		// 设置交易源数据Hash
		transactionCopy.SetDataHash(quotedTransactions, i, &input)

		publicKey := utils.GetPublicKey(transaction.TXInputs[i].PublicKey)
		dataHash := transactionCopy.TXId
		r, s := utils.GetRS(transaction.TXInputs[i].Signature)

		if !ecdsa.Verify(publicKey, dataHash, r, s) {
			fmt.Printf("交易验证失败！\n")
			return false
		}
	}

	return true
}

// 拷贝修剪的交易
func (transaction *Transaction) TrimmedCopy() *Transaction {
	var inputs []TXInput

	for _, input := range transaction.TXInputs {
		trimmedInput := TXInput{input.TXId, input.Index, []byte{}, []byte{}}
		inputs = append(inputs, trimmedInput)
	}

	transactionCopy := Transaction{transaction.TXId, inputs, transaction.TXOutputs}

	return &transactionCopy
}

// 设置交易源数据Hash（可用于交易签名和交易验证）
func (transaction *Transaction) SetDataHash(quotedTransactions map[string]*Transaction, i int, input *TXInput) {
	// 获取input引用的output所属的交易
	quotedTransaction := quotedTransactions[string(input.TXId)]
	// 获取input引用的output的PublicKeyHash赋值给transactionCopy对应input的PublicKey
	transaction.TXInputs[i].PublicKey = quotedTransaction.TXOutputs[input.Index].PublicKeyHash

	// 进行Hash运算
	transaction.SetTXId()

	// 重置transactionCopy对应input的PublicKey
	transaction.TXInputs[i].PublicKey = []byte{}
}

// 打印交易信息
func (transaction *Transaction) ToString() string {
	var lines []string

	for i, input := range transaction.TXInputs {
		lines = append(lines, fmt.Sprintf("  Input  %d", i))
		lines = append(lines, fmt.Sprintf("    TXId           %x", input.TXId))
		lines = append(lines, fmt.Sprintf("    Output index   %d", input.Index))
		lines = append(lines, fmt.Sprintf("    Signature      %x", input.Signature))
		lines = append(lines, fmt.Sprintf("    PublicKey      %x", input.PublicKey))
	}

	for i, output := range transaction.TXOutputs {
		lines = append(lines, fmt.Sprintf("  Output  %d", i))
		lines = append(lines, fmt.Sprintf("    Amount         %f", output.Amount))
		lines = append(lines, fmt.Sprintf("    PublicKeyHash  %x", output.PublicKeyHash))
	}

	return strings.Join(lines, "\n")
}