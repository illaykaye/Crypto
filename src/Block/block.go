package Block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hash"
	"time"
)

//Data class (transaction data)
type DataS struct {
	Sender    string  `json:"sender:"`
	Recipinet string  `json:"recipinet:"`
	Quantity  float64 `json:"quantiy:"`
}

func NewData(sender string, recipinet string, quantity float64) DataS {
	return DataS{
		sender,
		recipinet,
		quantity,
	}
}

func (d DataS) EncodeData() []byte {
	dataJSON, err := json.Marshal(d)
	isErr(err)
	return dataJSON
}

func DecodeData(dataJSON []byte) DataS {
	var data DataS
	err := json.Unmarshal(dataJSON, &data)
	isErr(err)

	return data
}

//Block class
type Block struct {
	Index     int       `json:"index:"`
	Proof_num int       `json:"proof_num:"`
	Prev_hash []byte    `json:"prev_hash:"`
	Data      []DataS   `json:"data:"`
	Time      time.Time `json:"time:"`
}

func NewBlock(index int, proof_num int, prev_hash []byte, data []DataS, time time.Time) Block {
	return Block{
		index,
		proof_num,
		prev_hash,
		data,
		time,
	}
}

func (b Block) CalcHash() []byte {
	var hash hash.Hash = sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d %d %v %v %s", b.Index, b.Proof_num, b.Prev_hash, b.Data, b.Time.Format(time.Stamp))))

	return hash.Sum(nil)
}

func (b Block) EncodeBlock() []byte {
	blockJSON, err := json.Marshal(b)
	isErr(err)
	return blockJSON
}

func DecodeBlock(blockJSON []byte) Block {
	var block Block
	err := json.Unmarshal(blockJSON, &block)
	isErr(err)

	return block
}

func isErr(err error) {
	if err != nil {
		panic(err)
	}
}
