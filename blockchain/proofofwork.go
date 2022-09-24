package blockchain

import (
	"BTC_Simulate/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 定义工作量证明结构
type ProofOfWork struct {
	block  *Block   // 区块
	target *big.Int // 目标值
}

// 生成目标值
func (pow *ProofOfWork) SetTarget() {
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"

	tempTarget := big.Int{}
	tempTarget.SetString(targetStr, 16)

	pow.target = &tempTarget
}

// 创建工作量证明
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	pow.SetTarget()

	return &pow
}

// 获取区块信息
func GetBlockInfo(block *Block, nonce uint64) []byte {
	// 拼接区块数据
	temp := [][]byte{
		utils.Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		utils.Uint64ToByte(block.TimeStamp),
		utils.Uint64ToByte(block.Difficulty),
		utils.Uint64ToByte(nonce),
		block.Data,
	}

	blockInfo := bytes.Join(temp, []byte{})

	return blockInfo
}

// 寻找随机数
func (pow *ProofOfWork) Run() (uint64, []byte) {
	var nonce uint64
	var hash [32]byte

	for {
		// 获取区块信息
		blockInfo := GetBlockInfo(pow.block, nonce)

		// 进行Hash运算
		hash = sha256.Sum256(blockInfo)

		// 将Hash数组转换成big.Int
		hashInt := big.Int{}
		hashInt.SetBytes(hash[:])

		// 将当前的Hash值与目标值进行比较
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！Hash值：%x，随机数：%d\n", hash, nonce)
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}
