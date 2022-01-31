package main

import (
	"github.com/devgony/nomadcoin/cli"
	"github.com/devgony/nomadcoin/db"
)

func main() {
	defer db.Close()
	// blockchain.Blockchain()
	cli.Start()

}
