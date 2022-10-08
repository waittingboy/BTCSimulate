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
	addressByte := base58.Decode(address)

	// 计算公钥Hash
	publicKeyHash := addressByte[1 : len(addressByte)-4]

	return publicKeyHash
}

// 通过publicKeyByte得到publicKey
func GetPublicKey(publicKeyByte []byte) *ecdsa.PublicKey {
	X := big.Int{}
	Y := big.Int{}

	X.SetBytes(publicKeyByte[:len(publicKeyByte)/2])
	Y.SetBytes(publicKeyByte[len(publicKeyByte)/2:])

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
