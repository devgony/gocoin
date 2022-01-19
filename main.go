package main

import (
	"github.com/devgony/nomadcoin/explorer"
	"github.com/devgony/nomadcoin/rest"
)

func main() {
	go explorer.Start(3002)
	rest.Start(4000)
}
