package blockchain

import (
	"github.com/boltdb/bolt"
	"log"
)

// 定义区块链结构
type BlockChain struct {
	db *bolt.DB
}

const BlockchainDB = "blockchain.db"
const BlockBucket = "blockBucket"
const LastHashKey = "lastHashKey"

// 创建创世块
func GenesisBlock() *Block {
	return NewBlock("我们的未来是星辰和大海！", []byte{})
}

// 创建区块链
func NewBlockChain() *BlockChain {
	// 打开数据库
	db, err := bolt.Open(BlockchainDB, 0600, nil)
	if err != nil {
		log.Panic("NewBlockChain：打开数据库失败！")
	}
	defer db.Close()

	// 写入数据库
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockBucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(BlockBucket))
			if err != nil {
				log.Panic("NewBlockChain：Bucket创建失败！")
			}

			genesisBlock := GenesisBlock()

			bucket.Put(genesisBlock.CurHash, genesisBlock.Serialize())
			bucket.Put([]byte(LastHashKey), genesisBlock.CurHash)
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
	db, err := bolt.Open(BlockchainDB, 0600, nil)
	if err != nil {
		log.Panic("AddBlock：打开数据库失败！")
	}
	defer db.Close()

	// 写入数据库
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockBucket))
		if bucket == nil {
			log.Panic("AddBlock：请先创建区块链！")
		}

		lastHash := bucket.Get([]byte(LastHashKey))
		block := NewBlock(data, lastHash)

		bucket.Put(block.CurHash, block.Serialize())
		bucket.Put([]byte(LastHashKey), block.CurHash)

		return nil
	})

	if err != nil {
		log.Panic("AddBlock：写入数据库失败！")
	}
}
