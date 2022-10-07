package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

// 定义区块头结构
type BlockHead struct {
	Version    uint64 // 版本号
	PrevHash   []byte // 前区块Hash
	MerkelRoot []byte // 梅克尔根
	TimeStamp  uint64 // 时间戳
	Difficulty uint64 // 难度值
	Nonce      uint64 // 随机数
}

// 定义区块结构
type Block struct {
	BlockHead    *BlockHead     // 当前区块头
	CurHash      []byte         // 当前区块Hash
	Transactions []*Transaction // 区块交易数组
}

// 模拟梅克尔根：只是对交易数据做简单的拼接，不做二叉树处理
func (block *Block) SetMerkelRoot() {
	var info []byte

	// 遍历区块交易
	for _, transaction := range block.Transactions {
		// 拼接交易ID
		info = append(info, transaction.TXId...)
	}

	// 进行Hash运算
	hash := sha256.Sum256(info)

	block.BlockHead.MerkelRoot = hash[:]
}

// 创建区块
func NewBlock(transactions []*Transaction, prevHash []byte) *Block {
	blockHead := BlockHead{
		Version:    00,
		PrevHash:   prevHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
	}

	block := Block{
		BlockHead:    &blockHead,
		CurHash:      []byte{},
		Transactions: transactions,
	}

	block.SetMerkelRoot()

	// 创建工作量证明
	pow := NewProofOfWork(&blockHead)
	// 寻找随机数
	nonce, curHash := pow.Run()
	// 更新block数据
	block.BlockHead.Nonce = nonce
	block.CurHash = curHash

	return &block
}

// 区块序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer

	// 定义编码器
	encoder := gob.NewEncoder(&buffer)
	// 使用编码器进行编码
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("区块编码失败！")
	}

	return buffer.Bytes()
}

// 区块反序列化
func Deserialize(data []byte) *Block {
	var block *Block

	// 定义解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	// 使用解码器进行解码
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("区块解码失败！")
	}

	return block
}
