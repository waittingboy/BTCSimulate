package blockchain

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	db      *bolt.DB
	curHash []byte
}

func (blockchain *Blockchain) NewIterator() *BlockchainIterator {
	return &BlockchainIterator{
		blockchain.db,
		blockchain.lastHash,
	}
}

func (iterator *BlockchainIterator) Next() *Block {
	var block *Block

	// 查询数据库
	err := iterator.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("请先创建区块链！\n")
		}

		blockBytes := bucket.Get(iterator.curHash)
		block = Deserialize(blockBytes)

		iterator.curHash = block.BlockHead.PrevHash

		return nil
	})

	if err != nil {
		log.Panic("查询数据库失败！")
	}

	return block
}
