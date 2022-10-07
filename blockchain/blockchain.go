package blockchain

import (
	"BTC_Simulate/utils"
	"BTC_Simulate/wallet"
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// 定义区块链结构
type Blockchain struct {
	db       *bolt.DB
	lastHash []byte
}

const blockchainDB = "blockchain.db"
const blockBucket = "blockBucket"
const lastHashKey = "lastHashKey"

// 创建创世块
func GenesisBlock(miner string) *Block {
	transaction := NewCoinbase(miner, "我们的未来是星辰和大海！")
	return NewBlock([]*Transaction{transaction}, []byte{})
}

// 创建区块链
func NewBlockchain(miner string) *Blockchain {
	if !wallet.IsValidAddress(miner) {
		fmt.Printf("miner地址无效，请重新输入！\n")
		return nil
	}

	// 打开数据库
	db, err := bolt.Open(blockchainDB, 0600, nil)
	if err != nil {
		log.Panic("NewBlockchain：打开数据库失败！")
	}

	// 写入数据库
	var lastHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("NewBlockchain：Bucket创建失败！")
			}

			genesisBlock := GenesisBlock(miner)

			bucket.Put(genesisBlock.CurHash, genesisBlock.Serialize())
			bucket.Put([]byte(lastHashKey), genesisBlock.CurHash)

			lastHash = genesisBlock.CurHash
		} else {
			lastHash = bucket.Get([]byte(lastHashKey))
		}

		return nil
	})

	if err != nil {
		log.Panic("NewBlockchain：写入数据库失败！")
	}

	return &Blockchain{db, lastHash}
}

// 添加区块
func (blockchain *Blockchain) AddBlock(transactions []*Transaction) {
	// 写入数据库
	err := blockchain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("请先创建区块链！\n")
		}

		lastHash := bucket.Get([]byte(lastHashKey))
		block := NewBlock(transactions, lastHash)

		bucket.Put(block.CurHash, block.Serialize())
		bucket.Put([]byte(lastHashKey), block.CurHash)

		return nil
	})

	if err != nil {
		log.Panic("AddBlock：写入数据库失败！")
	}
}

// 查找包含用户所有UTXO的交易集合
func (blockchain *Blockchain) FindUTXOTransaction(user string) []*Transaction {
	// 包含用户UTXO的交易集合
	var transactions []*Transaction
	// 已消耗的output集合
	spentOutputs := make(map[string][]int)

	iterator := blockchain.NewIterator()
	for {
		block := iterator.Next()
		// 一个用户在一个交易中只会存在input或output中的一个，不会同时存在input和output
		// 遍历区块交易：可筛选出所有包含当前用户UTXO的交易
		for _, transaction := range block.Transactions {
			spentIndexes := spentOutputs[string(transaction.TXId)]
			// 遍历outputs：必须先遍历outputs再遍历inputs，最后一个区块的所有outputs都是UTXO
			for i, output := range transaction.TXOutputs {
				// 将包含用户UTXO的交易加入transactions
				if bytes.Equal(output.PublicKeyHash, GetPublicKeyHash(user)) {
					if !utils.Contains(i, spentIndexes) {
						transactions = append(transactions, transaction)
					}
				}
			}

			// 遍历输入：不遍历Coinbase交易的输入
			if !transaction.IsCoinbase() {
				for _, input := range transaction.TXInputs {
					if bytes.Equal(wallet.GetPublicKeyHash(input.PublicKey), GetPublicKeyHash(user)) {
						spentOutputs[string(input.TXId)] = append(spentOutputs[string(input.TXId)], input.Index)
					}
				}
			}
		}

		if len(block.BlockHead.PrevHash) == 0 {
			break
		}
	}

	return transactions
}

// 查找用户所有UTXO
func (blockchain *Blockchain) FindUTXOs(user string) []TXOutput {
	// 所有UTXO集合
	var UTXOs []TXOutput

	// 包含用户所有UTXO的交易集合
	transactions := blockchain.FindUTXOTransaction(user)

	// 遍历交易集合
	for _, transaction := range transactions {
		// 遍历outputs，将用户的UTXO加入UTXOs
		for _, output := range transaction.TXOutputs {
			if bytes.Equal(output.PublicKeyHash, GetPublicKeyHash(user)) {
				UTXOs = append(UTXOs, output)
			}
		}
	}

	return UTXOs
}

// 查找用户转账所需的最少UTXOs
func (blockchain *Blockchain) FindNeedUTXOs(user string, amount float64) (map[string][]int, float64) {
	// 所需的最少UTXO集合
	UTXOs := make(map[string][]int)

	// 最少UTXO集合的总金额
	var totalAmount float64

	// 包含用户所有UTXO的交易集合
	transactions := blockchain.FindUTXOTransaction(user)

	// 遍历交易集合
	for _, transaction := range transactions {
		// 遍历outputs，将用户所需的UTXO加入UTXOs
		for i, output := range transaction.TXOutputs {
			if bytes.Equal(output.PublicKeyHash, GetPublicKeyHash(user)) {
				UTXOs[string(transaction.TXId)] = append(UTXOs[string(transaction.TXId)], i)
				totalAmount += output.Amount
				if totalAmount >= amount {
					return UTXOs, totalAmount
				}
			}
		}
	}

	return UTXOs, totalAmount
}
