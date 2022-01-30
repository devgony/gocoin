package main

import "github.com/devgony/nomadcoin/blockchain"

func main() {
	// cli.Start()
	blockchain.Blockchain().AddBlock("First")
	blockchain.Blockchain().AddBlock("Second")
	blockchain.Blockchain().AddBlock("Third")
}
