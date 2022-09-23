package blockchain

import (
	"BTC_Simulate/utils"
	"bytes"
	"crypto/sha256"
	"time"
)

// 定义区块结构
type Block struct {
	Version    uint64 // 版本号
	PrevHash   []byte // 前区块Hash
	MerkelRoot []byte // 梅克尔根
	TimeStamp  uint64 // 时间戳
	Difficulty uint64 // 难度值
	Nonce      uint64 // 随机数
	CurHash    []byte // 当前区块Hash
	Data       []byte // 区块数据
}

// 生成当前区块Hash(sha256)
func (block *Block) setHash() {
	// 拼接区块数据
	temp := [][]byte{
		utils.Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		utils.Uint64ToByte(block.TimeStamp),
		utils.Uint64ToByte(block.Difficulty),
		block.Data,
	}

	blockInfo := bytes.Join(temp, []byte{})

	// 生成区块Hash
	hash := sha256.Sum256(blockInfo)
	block.CurHash = hash[:]
}

// 创建区块
func NewBlock(data string, prevHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		CurHash:    []byte{},
		Data:       []byte(data),
	}

	block.setHash()

	return &block
}
