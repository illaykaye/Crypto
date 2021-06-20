package Block

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"time"
)

type BlockChain struct {
	currentData []DataS
	file        string
	lastBlock   Block
	numBlocks   int
}

func NewBC(file string) BlockChain {
	var (
		numLines int   = numLinesInFile(file)
		block    Block = Block{}
	)
	fmt.Println(numLines)
	if numLines > 0 {
		//fmt.Println(string(readLine(file, numLines)))
		block = DecodeBlock(readLine(file, numLines))
	}

	return BlockChain{
		make([]DataS, 0),
		file,
		block,
		numLines,
	}
}

func (bc *BlockChain) GetNumBlocks() int {
	return bc.numBlocks
}

func (bc *BlockChain) GetCurrentData() []DataS {
	return bc.currentData
}

func (bc *BlockChain) GetLatestBlock() Block {
	return bc.lastBlock
}

func (bc *BlockChain) InitLastBlock() {
	blockJSON := readLine(bc.file, bc.numBlocks+1)
	bc.lastBlock = DecodeBlock(blockJSON)
}

func readLine(f string, lineNum int) []byte {
	var lastLine int = 1
	file, err := os.OpenFile(f, os.O_RDONLY, 0644)
	isErr(err)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		if lastLine == lineNum {
			//fmt.Println(sc.Text())
			return sc.Bytes()
		}
		lastLine++
	}
	return []byte{}
}

func (bc *BlockChain) openFile() *os.File {
	f, err := os.OpenFile(bc.file, os.O_APPEND|os.O_RDONLY, 0644)
	isErr(err)
	return f
}

func (bc *BlockChain) AddData(sender string, recipinet string, quantity float64) {
	var data DataS = NewData(sender, recipinet, quantity)

	bc.currentData = append(bc.currentData, data)
}

func (bc *BlockChain) ConstructBlock() Block {
	var block Block = NewBlock(bc.numBlocks, bc.newProofOfWork(), bc.GetLatestBlock().CalcHash(), bc.currentData, time.Now())
	fmt.Println("rawr")
	//bc.chain = append(bc.chain, block)
	bc.currentData = make([]DataS, 0)

	bc.lastBlock = block
	bc.writeLatestBlock()

	return block
}

func (bc *BlockChain) Genesis() Block {
	var hash hash.Hash = sha256.New()
	hash.Write([]byte{})

	var block Block = NewBlock(0, 0, []byte("0"), nil, time.Now())

	//bc.chain = append(bc.chain, block)
	bc.lastBlock = block
	bc.numBlocks++
	bc.writeLatestBlock()
	return block
}

func (bc *BlockChain) newProofOfWork() int {
	var proof_no int = 0
	var last_proof int = bc.GetLatestBlock().Proof_num

	for !verifyingProof(proof_no, last_proof) {
		proof_no++
	}

	return proof_no
}

func verifyingProof(lastProof int, proof int) bool {
	var hash hash.Hash = sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d %d", lastProof, proof)))

	bs := hash.Sum(nil)
	var hex string = hex.EncodeToString(bs)

	return hex[0:4] == "0000"
}

func CheckValidity(block Block, prevBlock Block) bool {
	if prevBlock.Index+1 != block.Index {
		return false
	} else if string(block.Prev_hash) != string(prevBlock.CalcHash()) {
		return false
	} else if block.Time.Before(prevBlock.Time) {
		return false
	} else if !verifyingProof(block.Proof_num, prevBlock.Proof_num) {
		return false
	}

	return true
}

/*
func (bc *BlockChain) InitBlockChainFromFile() {
	var (
		reader   *bufio.Reader = bufio.NewReader(bc.file)
		lineNum  int           = 0
		numLines int           = numLinesInFile(bc.file)
		block    Block         = Block{}
	)

	for lineNum < numLines {
		blockJSON, _, err := reader.ReadLine()
		isErr(err)

		block = DecodeBlock(blockJSON)
		bc.AddBlock(block)

		lineNum++
	}

}
*/

func (bc *BlockChain) writeLatestBlock() {
	f := bc.openFile()
	f.WriteString(fmt.Sprintf("%s\n", string(bc.GetLatestBlock().EncodeBlock())))
	f.Close()
}

func numLinesInFile(file string) int {
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	isErr(err)

	fileScanner := bufio.NewScanner(f)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	f.Close()
	return lineCount
}
