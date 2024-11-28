package hashing

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/wealdtech/go-merkletree/keccak256"
)

func ToSha256(password string) string {
	sha256Hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sha256Hash[:])
}

func ToKeccak256(password string) string {
	k256 := keccak256.Keccak256{}
	k256Hash := k256.Hash([]byte(password))
	return hex.EncodeToString(k256Hash[:])
}

func ToMD5(password string) string {
	md5Hash := md5.Sum([]byte(password))
	return hex.EncodeToString(md5Hash[:])
}
