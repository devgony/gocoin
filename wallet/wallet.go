package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/devgony/nomadcoin/utils"
)

const (
	fileName string = "nomadcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat("nomadcoin.wallet")
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

func aFromK(key *ecdsa.PrivateKey) string {
	z := append(key.X.Bytes(), key.Y.Bytes()...)
	return fmt.Sprintf("%x", z)
}

func sign(payload string, w *wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	signature := append(r.Bytes(), s.Bytes()...)
	return fmt.Sprintf("%x", signature)
}

func verify(signature string, payload string, publicKey string) bool {
	return true
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
