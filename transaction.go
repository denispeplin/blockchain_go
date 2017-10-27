package main

import (
	"crypto/sha256"
	"fmt"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	VOut []TXOutput
}

func (tx *Transaction) SetID() {
	var hash [32]byte

	encoded := Serialize(tx)

	hash = sha256.Sum256(encoded)

	tx.ID = hash[:]
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Revard to %s", to)
	}

	txin := TXInput{Txid: []byte{}, Vout: -1, ScriptSig: data}
	txout := TXOutput{Value: subsidy, ScriptPubKey: to}
	tx := Transaction{ID: nil, Vin: []TXInput{txin}, VOut: []TXOutput{txout}}

	tx.SetID()

	return &tx
}
