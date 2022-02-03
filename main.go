package main

import (
	"github.com/devgony/nomadcoin/cli"
	"github.com/devgony/nomadcoin/db"
)

const (
	privateKey    string = "30770201010420270623da3768df6fc3c3439b8e0319621318b1dec6199052f49faefdd9d80548a00a06082a8648ce3d030107a1440342000462ded99b11da850eec19a908aa57effbec88541aa04da07d0a2cabf046b2502dd061eccc9860c7922ea758a2e8ac1e5f6d044d7a6af03060aa5dcb13cafc8a73"
	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
	signature     string = "6d56582490ff9a54b44df6bf9fa991c0432f2fd25f32760bf540b10049b50a048ca986d7e9f0ee7745bce735dcd0db951f21664f054e94b0a03d87046857a3ca"
)

func main() {
	defer db.Close()
	// blockchain.Blockchain()
	cli.Start()
	// wallet.Wallet()
}
