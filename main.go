package main

import (
	"github.com/devgony/nomadcoin/blockchain"
	"github.com/devgony/nomadcoin/cli"
	"github.com/devgony/nomadcoin/db"
)

func main() {
	defer db.Close()
	blockchain.Blockchain()
	cli.Start()
	// blockchain.Blockchain().AddBlock("First")
	// blockchain.Blockchain().AddBlock("Second")
	// blockchain.Blockchain().AddBlock("Third")
}
