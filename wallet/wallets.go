package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.dat"

// 定义钱包集合结构
type Wallets struct {
	WalletsMap map[string]*Wallet
}

// 创建钱包集合
func NewWallets() *Wallets {
	var wallets Wallets

	wallets.loadFile()

	return &wallets
}

// 添加钱包
func (wallets *Wallets) AddWallet() string {
	wallet := NewWallet()
	address := wallet.GenerateAddress()
	wallets.WalletsMap[address] = wallet

	wallets.saveToFile()

	return address
}

// 加载钱包集合
func (wallets *Wallets) loadFile() {
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		wallets.WalletsMap = make(map[string]*Wallet)
		return
	}

	walletsBytes, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic("钱包集合数据读取失败！\n")
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(walletsBytes))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic("钱包集合解码失败！\n")
	}
}

// 保存钱包集合
func (wallets *Wallets) saveToFile() {
	var buffer bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic("钱包集合编码失败！\n")
	}

	err = ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
	if err != nil {
		log.Panic("钱包集合数据写入失败！\n")
	}
}

// 列出所有地址
func (wallets *Wallets) ListAllAddresses() []string {
	var addresses []string

	for address, _ := range wallets.WalletsMap {
		addresses = append(addresses, address)
	}

	return addresses
}
