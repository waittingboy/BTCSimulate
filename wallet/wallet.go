package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

// 定义钱包结构
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey // 私钥
	PublicKey  *ecdsa.PublicKey  // 公钥
}

// 创建钱包
func NewWallet() *Wallet {
	// 创建私钥
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic("生成私钥失败！\n")
	}

	// 创建公钥
	publicKey := privateKey.PublicKey

	// 创建钱包
	wallet := Wallet{privateKey, &publicKey}

	return &wallet
}

// 生成地址
func (wallet *Wallet) GenerateAddress() string {
	// 获取公钥Hash
	publicKey := wallet.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	publicKeyHash := GetPublicKeyHash(publicKeyBytes)

	// 拼接version与publicKeyHash
	version := []byte{00}
	payload := append(version, publicKeyHash...)

	// 获取校验码
	checkSum := GetCheckSum(payload)

	// 拼接payload与校验码
	addressBytes := append(payload, checkSum...)

	// 进行Base58转换得到地址
	address := base58.Encode(addressBytes)

	return address
}

// 计算公钥Hash
func GetPublicKeyHash(data []byte) []byte {
	// 进行hash运算
	hash := sha256.Sum256(data)

	// 进行RIPEMD160运算
	rip160 := ripemd160.New()
	_, err := rip160.Write(hash[:])
	if err != nil {
		log.Panic("RIPEMD160运算失败!\n")
	}
	publicKeyHash := rip160.Sum([]byte{})

	return publicKeyHash
}

// 计算校验码
func GetCheckSum(data []byte) []byte {
	// 对data进行双hash运算
	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])

	// 取前四个字节作为校验码
	checkSum := hash[:4]

	return checkSum
}

// 地址校验
func IsValidAddress(address string) bool {
	// 解码
	addressBytes := base58.Decode(address)
	if len(addressBytes) < 4 {
		return false
	}

	// 拆分数据
	payload := addressBytes[:len(addressBytes)-4]
	checkSum := addressBytes[len(addressBytes)-4:]

	// 计算checkSum并比较
	return bytes.Equal(GetCheckSum(payload), checkSum)
}
