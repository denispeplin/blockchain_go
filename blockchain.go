package main

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	newBlock := NewBlock(transactions, bc.tip)

	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		RaiseError(err)
		err = bucket.Put([]byte("l"), newBlock.Hash)
		RaiseError(err)
		bc.tip = newBlock.Hash

		return nil
	})
	RaiseError(err)
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := BlockchainIterator{currentHash: bc.tip, db: bc.db}

	return &bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	RaiseError(err)

	i.currentHash = block.PrevBlockHash

	return block
}

func dbExists() bool {
	_, err := os.Stat(dbFile)
	return os.IsExist(err)
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	RaiseError(err)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		tip = bucket.Get([]byte("l"))

		return nil
	})
	RaiseError(err)

	return &Blockchain{tip: tip, db: db}
}

func CreateBlockChain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	RaiseError(err)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			//fmt.Println("No existing blockchain found. Creating a new one...")
			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)
			bucket, err = tx.CreateBucket([]byte(blocksBucket))
			RaiseError(err)
			err = bucket.Put(genesis.Hash, genesis.Serialize())
			RaiseError(err)
			err = bucket.Put([]byte("l"), genesis.Hash)
			RaiseError(err)
			tip = genesis.Hash
		}

		return nil
	})
	RaiseError(err)

	return &Blockchain{tip: tip, db: db}
}
