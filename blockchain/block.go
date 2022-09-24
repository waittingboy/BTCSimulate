package blockchain

import "time"

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

	// 创建工作量证明
	pow := NewProofOfWork(&block)
	// 寻找随机数
	nonce, curHash := pow.Run()
	// 更新block数据
	block.Nonce = nonce
	block.CurHash = curHash

	return &block
}
