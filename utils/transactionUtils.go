package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/btcsuite/btcd/btcutil/base58"
	"math/big"
)

// 计算公钥Hash
func GetPublicKeyHash(address string) []byte {
	// 解码
	addressBytes := base58.Decode(address)

	// 计算公钥Hash
	publicKeyHash := addressBytes[1 : len(addressBytes)-4]

	return publicKeyHash
}

// 通过publicKeyByte得到publicKey
func GetPublicKey(publicKeyBytes []byte) *ecdsa.PublicKey {
	X := big.Int{}
	Y := big.Int{}

	X.SetBytes(publicKeyBytes[:len(publicKeyBytes)/2])
	Y.SetBytes(publicKeyBytes[len(publicKeyBytes)/2:])

	curve := elliptic.P256()
	publicKey := ecdsa.PublicKey{Curve: curve, X: &X, Y: &Y}

	return &publicKey
}

// 通过数字签名拆分得到r,s
func GetRS(signature []byte) (*big.Int, *big.Int) {
	r := big.Int{}
	s := big.Int{}

	r.SetBytes(signature[:len(signature)/2])
	s.SetBytes(signature[len(signature)/2:])

	return &r, &s
}
