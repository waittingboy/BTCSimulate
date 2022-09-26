package blockchain

import (
	"github.com/boltdb/bolt"
	"log"
)

// 定义区块链结构
type BlockChain struct {
	db *bolt.DB
}

const blockchainDB = "blockchain.db"
const blockBucket = "blockBucket"
const lastHashKey = "lastHashKey"

// 创建创世块
func GenesisBlock() *Block {
	return NewBlock("我们的未来是星辰和大海！", []byte{})
}

// 创建区块链
func NewBlockChain() *BlockChain {
	// 打开数据库
	db, err := bolt.Open(blockchainDB, 0600, nil)
	if err != nil {
		log.Panic("NewBlockChain：打开数据库失败！")
	}
	defer db.Close()

	// 写入数据库
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("NewBlockChain：Bucket创建失败！")
			}

			genesisBlock := GenesisBlock()

			bucket.Put(genesisBlock.CurHash, genesisBlock.Serialize())
			bucket.Put([]byte(lastHashKey), genesisBlock.CurHash)
		}

		return nil
	})

	if err != nil {
		log.Panic("NewBlockChain：写入数据库失败！")
	}

	return &BlockChain{db}
}

// 添加区块
func (blockChain *BlockChain) AddBlock(data string) {
	// 打开数据库
	db, err := bolt.Open(blockchainDB, 0600, nil)
	if err != nil {
		log.Panic("AddBlock：打开数据库失败！")
	}
	defer db.Close()

	// 写入数据库
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("AddBlock：请先创建区块链！")
		}

		lastHash := bucket.Get([]byte(lastHashKey))
		block := NewBlock(data, lastHash)

		bucket.Put(block.CurHash, block.Serialize())
		bucket.Put([]byte(lastHashKey), block.CurHash)

		return nil
	})

	if err != nil {
		log.Panic("AddBlock：写入数据库失败！")
	}
}
