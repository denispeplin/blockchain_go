package main

func main() {
	bc := NewBlockchain()
	defer bc.db.Close()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 1 more BTC to Ivan")
}
