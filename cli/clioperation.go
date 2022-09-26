package cli

import (
	"BTC_Simulate/blockchain"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func (cli *CLI) AddBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Printf("添加区块成功！")
}

func (cli *CLI) PrintBlockchain() {
	// 打开数据库
	db, err := bolt.Open(blockchain.BlockchainDB, 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败！")
	}
	defer db.Close()

	// 查询数据库
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockchain.BlockBucket))
		if bucket == nil {
			log.Panic("请先创建区块链！")
		}

		lastHash := bucket.Get([]byte(blockchain.LastHashKey))
		if len(lastHash) != 0 {
			count := 0

			for {
				count++
				value := bucket.Get(lastHash)
				block := blockchain.Deserialize(value)
				fmt.Printf("========================================\n")
				fmt.Printf("当前版本号：%d\n", block.Version)
				fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
				fmt.Printf("区块时间戳：%d\n", block.TimeStamp)
				fmt.Printf("随机数：%d\n", block.Nonce)
				fmt.Printf("当前区块哈希值：%x\n", block.CurHash)
				fmt.Printf("区块数据：%s\n", block.Data)

				if len(block.PrevHash) == 0 {
					break
				}

				lastHash = block.PrevHash
			}
		}

		return nil
	})

	if err != nil {
		log.Panic("查询数据库失败！")
	}
}