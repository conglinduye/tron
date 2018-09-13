package utils

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/core"
)

// CalcTransactionHash 计算Transaction的hash, 返回为[]byte, 显示时需要使用 hex encoding 编码
func CalcTransactionHash(transaction *core.Transaction) []byte {
	if nil == transaction {
		return nil
	}
	rawByte, _ := proto.Marshal(transaction.RawData)
	hash0 := sha256.Sum256(rawByte)
	return hash0[:]
}

// CalcBlockHash 计算Block的Hash, 返回为[]byte, 显示时需要使用 hex encoding 编码
func CalcBlockHash(block *core.Block) (ret []byte) {
	if nil == block || nil == block.BlockHeader {
		return nil
	}
	rawByte, _ := proto.Marshal(block.BlockHeader.RawData) // get block header raw bytes
	hash0 := sha256.Sum256(rawByte)                        // calc hash using sha256
	// binary.BigEndian.PutUint64(numByte, uint64(block.BlockHeader.RawData.Number)) // convert block num to byte
	// copy(hash0[0:8], numByte)                                                     // replace first 8 byte by num byte
	binary.BigEndian.PutUint64(hash0[:], uint64(block.BlockHeader.RawData.Number))
	return hash0[:]
}

// CalcBlockSize 计算Block的size
func CalcBlockSize(block *core.Block) int64 {
	if nil == block || nil == block.BlockHeader || nil == block.BlockHeader.RawData {
		return 0
	}
	rawByte, _ := proto.Marshal(block) // get block header raw bytes
	return int64(len(rawByte))
}
