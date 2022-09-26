package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
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

// 区块序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer

	// 定义编码器
	encode := gob.NewEncoder(&buffer)
	// 使用编码器进行编码
	err := encode.Encode(&block)
	if err != nil {
		log.Panic("区块编码失败！")
	}

	return buffer.Bytes()
}

// 区块反序列化
func Deserialize(data []byte) *Block {
	var block *Block

	// 定义解码器
	decode := gob.NewDecoder(bytes.NewReader(data))
	// 使用解码器进行解码
	err := decode.Decode(&block)
	if err != nil {
		log.Panic("区块解码失败！")
	}

	return block
}
