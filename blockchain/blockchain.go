package blockchain

// 定义区块链结构
type BlockChain struct {
	Blocks []*Block // 区块数组
}

// 创建创世块
func GenesisBlock() *Block {
	return NewBlock("我们的未来是星辰和大海！", []byte{})
}

// 创建区块链
func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()

	blockChain := BlockChain{
		Blocks: []*Block{genesisBlock}, // 添加创世块
	}

	return &blockChain
}

// 添加区块
func (blockChain *BlockChain) AddBlock(data string) {
	// 获取前区块Hash
	lastBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
	prevHash := lastBlock.CurHash

	// 创建新区块
	block := NewBlock(data, prevHash)

	// 添加到区块数组中
	blockChain.Blocks = append(blockChain.Blocks, block)
}
