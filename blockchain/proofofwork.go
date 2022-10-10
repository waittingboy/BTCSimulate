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
	blockHead *BlockHead // 区块头
	target    *big.Int   // 目标值
}

// 生成目标值
func (pow *ProofOfWork) SetTarget() {
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"

	tempTarget := big.Int{}
	tempTarget.SetString(targetStr, 16)

	pow.target = &tempTarget
}

// 创建工作量证明
func NewProofOfWork(blockHead *BlockHead) *ProofOfWork {
	pow := ProofOfWork{
		blockHead: blockHead,
	}
	pow.SetTarget()

	return &pow
}

// 获取区块头信息
func GetBlockHeadInfo(blockHead *BlockHead, nonce uint64) []byte {
	// 拼接区块头数据
	temp := [][]byte{
		utils.Uint64ToBytes(blockHead.Version),
		blockHead.PrevHash,
		blockHead.MerkelRoot,
		utils.Uint64ToBytes(blockHead.TimeStamp),
		utils.Uint64ToBytes(blockHead.Difficulty),
		utils.Uint64ToBytes(nonce),
	}

	blockInfo := bytes.Join(temp, []byte{})

	return blockInfo
}

// 寻找随机数
func (pow *ProofOfWork) Run() (uint64, []byte) {
	var nonce uint64
	var hash [32]byte

	fmt.Printf("开始挖矿...\n")
	for {
		// 获取区块头信息
		blockHeadInfo := GetBlockHeadInfo(pow.blockHead, nonce)

		// 进行Hash运算
		hash = sha256.Sum256(blockHeadInfo)

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
