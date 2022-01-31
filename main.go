package main

import (
	"github.com/devgony/nomadcoin/blockchain"
	"github.com/devgony/nomadcoin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
	// blockchain.Blockchain().AddBlock("First")
	// blockchain.Blockchain().AddBlock("Second")
	// blockchain.Blockchain().AddBlock("Third")
}
