# 2. INTRODUCTION

## 2.2 Why Go (06:17)

- Only 1 option: for loop
- Powerful standard library
- Low learning curve

# 3. TOUR OF GO

## 3.0 Creating the Project (04:30)

```sh
go mod init github.com/devgony/nomadcoin
ls
> go.mod # like package.json
touch main.go
```

## 3.1 Variables in Go (07:54)

### create and update var syntax (should be inside func only)

```go
var name string = "henry"
name := "henry" // same with above, syntax sugar
```

- var, const
- bool, string, int(8,16..64), uint(8,16..64) byte, float(32,64)

## 3.2 Functions (08:59)

- If params have same type, specify only at the last
- func can return two types

```go
func plus(a, b int, name string) (int, string) {
	return a + b, name
}
```

- multiple params

```go
func plus(a ...int) int {
	var total int
	for index, item := range a {
		total += item
	}
	return total
}
```

## 3.3 fmt (03:52)

```go
x := 84375983402
fmt.Printf("%b\n", x)
fmt.Printf(fmt.Sprintf("%b\n", x)) // return fmted string (not print)
fmt.Printf("%o\n", x)
fmt.Printf("%x\n", x)
fmt.Printf("%U\n", x)
```

## 3.4 Slices and Arrays (08:02)

- array is declarative and limited in go

```go
foods := [3]string{"p", "o", "s"}
for i := 0; i < len(foods); i++ {
    fmt.Println(foods[i])
}
```

- slice is growable and infinited

```go
foods := []string{"p", "o", "s"}
fmt.Printf("%v\n", foods)
foods = append(foods, "t") // returns appended slice (should set to var manually)
fmt.Printf("%v\n", foods)
```

## 3.5 Pointers (08:44)

```go
a := 2
b := a  // copy
c := &a // borrow
a = 9
fmt.Println(a, b, *c) // 9 2 9
```

## 3.6 Structs (07:11)

### struct (like class)

```go
type person struct {
	name string
	age  int
}
```

### receiver function (like methods)

- choose first letter of type as a alias

```go
func (p person) sayHello() {
	// give access to instance
	fmt.Printf("Hello! My name is %s and I'm %d", p.name, p.age)
}
```

## 3.7 Structs with Pointers (09:35)

```sh
mkdir person
touch person/person.go
```

```go
package person

type Person struct {
	name string
	age  int
}

func (p Person) SetDetails(name string, age int) {
    // p is copy
	p.name = name
	p.age = age
}
```

### Receiver pointer function

- basically receiver function's p is copy -> use only for read
- we don't want to mutate copied, but origin
- use \* to mutate original

```go
func (p *Person) SetDetails(name string, age int) {
    // p is the origin
	p.name = name
	p.age = age
```

# 4. BLOCKCHAIN

## 4.0 Introduction (05:05)

- Concentrate to blockchain concept, solve side problem later

## 4.1 Our First Block (13:58)

- if any block is edited, invalid

```
b1Hash = (data + "")
b2Hash = (data + b1Hash)
b3Hash = (data + b2Hash)
```

- sha256 needs slice of bytes: cuz string is immutable

```go
genesisBlock := block{"Genesis Block", "", ""}
hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevHash))
hexHash := fmt.Sprintf("%x", hash)
genesisBlock.hash = hexHash
secondBlocks := block{"Second Blocks", "", genesisBlock.hash}
```

## 4.2 Our First Blockchain (09:43)

```go
type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

func (b *blockchain) getLastHash() string {
	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1].hash
	}
	return ""
}

func (b *blockchain) addBlock(data string) {
	newBlock := block{data, "", b.getLastHash()}
	hash := sha256.Sum256([]byte(newBlock.data + newBlock.prevHash))
	newBlock.hash = fmt.Sprintf("%x", hash)
	b.blocks = append(b.blocks, newBlock)
}

func (b *blockchain) listBlocks() {
	for _, block := range b.blocks {
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("PrevHash: %s\n", block.prevHash)
	}
}

func main() {
	chain := blockchain{}
	chain.addBlock("Genesis Block")
	chain.addBlock("Second Block")
	chain.addBlock("Third Block")
	chain.listBlocks()
}
```

## 4.3 Singleton Pattern (05:57)

```sh
mkdir blockchain
touch blockchain/blockchain.go
```

- singletone: share only 1 instance

```go
// blockchain/blockchain.go
package blockchain

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

var b *blockchain

func GetBlockchain() *blockchain {
	if b == nil {
		b = &blockchain{}
	}
	return b
}
```

## 4.4 Refactoring part One (09:16)

- Package sync.once: keep running once though ran by goroutine

```go
once.Do(func() {
			b = &blockchain{}
			b.blocks = append(b.blocks, createBlock(("Genesis Block")))
})
```

- Blockchain should be a slice of pointer with borrow (it will be way longer)

## 4.5 Refactoring part Two (07:15)

```go
package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*block
}

var b *blockchain
var once sync.Once

func (b *block) calculateHash() {
	Hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", Hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(Data string) *block {
	newBlock := block{Data, "", getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock((data)))

}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock(("Genesis"))
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*block {
	return b.blocks
}
```

# 5. EXPLORER

## 5.0 Setup (06:42)

- server side rendering only with std lib

```go
const port string = ":4000"

func home(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hello from home!") // print to writer
}

func main() {
	http.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil)) // if failed exit 1 else none
}
```

## 5.1 Rendering Templates (08:10)

```sh
mkdir templates
touch templates/home.html
```

- template.Must

```go
tmpl, err := template.ParseFiles("templates/home.html")
if err != nil {
		log.Fatal((err))
}
```

```go
tmpl := template.Must(template.ParseFiles("templates/home.html"))
```

- template

```go
type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	tmpl.Execute(rw, data)
}
```

## 5.2 Rendering Blocks (07:09)

- install extension: `gotemplate-syntax`
- mvp.css: https://andybrewer.github.io/mvp/

```html
<link rel="stylesheet" href="https://unpkg.com/mvp.css" />
```

- just copy & paste html? => partials => glob import

```sh
mkdir templates/partials
touch templates/partials/footer.gohtml
touch templates/partials/head.gohtml
touch templates/partials/header.gohtml
touch templates/partials/block.gohtml

mkdir templates/pages
move templates/home.gohtml templates/pages/home.gohtml
touch templates/pages/add.gohtml
```

## 5.3 Using Partials (10:34)

```sh
touch partials/block.gohtml
```

- load

```go
{{template "head"}}
```

- load glob
  - can't use `**/` so that parseGlob n times

```go
templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml")) // template^s
```

## 5.4 Adding A Block (14:44)

### Cursor synatax `.`: passing struct

1. template of template use `.`

```go
// home.gohtml
{{template "header" .PageTitle}}
```

```go
// header.gohtml
...
<h1>{{.}}</h1>
...
```

2. inside loop, use `.`

```go
{{range .Blocks}}
	{{template "block" .}}
{{end}}
```

- switch case "GET", "POST"

```go
func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}
```

## 5.5 Refactoring (04:42)

```sh
mkdir explorer
mv templates explorer/
cp main.go explorer/explorer.go
```

```go
// explorer/explorer.go
	templateDir string = "explorer/templates/"
```

# 6. REST API

## 6.0 Setup (09:03)

- REST API

```go
mkdir utils
touch utils/utils.go

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
```

### `json.Marshal(data)`

- Marshal: convert from goInterface to JSON

## 6.1 Marshal and Field Tags (11:18)

- manual marshal

```go
rw.Header().Add("Content-Type", "application/json")
b, err := json.Marshal(data)
utils.HandleErr(err)
fmt.Fprintf(rw, "%s", b)
```

- simple marshal

```go
json.NewEncoder(rw).Encode(data)
```

- struct field tag  
  https://pkg.go.dev/encoding/json#Marshal

```go
Description string `json:"description"` // make lowercase
Payload     string `json:"payload,omitempty"` // omit if empty
...
Payload:     "data:string", // write data on body
```

## 6.2 MarshalText (13:38)

### Interface

- if impl Stringer, can control fmt

```go
func (u URLDescription) String() string {
	return "Hello I'm a URL Description"
}
```

- prepend `http://localhost`

```go
func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}
```

## 6.3 JSON Decode (14:00)

- Install VSC extension: REST client

```sh
## touch api.http
http://localhost:4000/blocks
## send request by clicking
POST http://localhost:4000/blocks
{
    "message": "Data for my block"
}
```

- should pass pointer(address) to decode

```go
utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
```

## 6.4 NewServeMux (11:50)

- quick refactoring

```sh
mkdir rest
sed 's/main()/Start()/; s/package main/package rest/' main.go > rest/rest.go
echo "package main\nfunc main() {}" > main.go
```

- dynamic port

```go
var port string
...
func Start(aPort int) {
	port = fmt.Sprint(":%d", aPort)
```

```go
// explorer.go
fmt.Printf("Listening on http://localhost:%d\n", port)
log.Fatal(http.ListenAndServe(fmt.Sprint(":%d", port), nil)
```

### use NewServeMux

- to solve duped route
- nil -> defaultServeMux, handler -> NewServeMux

```go
// rest.go
handler := http.NewServeMux()
...
handler.HandleFunc("/", documentation)
handler.HandleFunc("/blocks", blocks)
...
log.Fatal(http.ListenAndServe(port, handler))
```

## 6.5 Gorilla Mux (08:54)

- can handle params

```sh
go get -u github.com/gorilla/mux
```

```go
router := mux.NewRouter()
...
vars := mux.Vars(r)
```

## 6.6 Atoi (08:42)

- string to int

```go
id, err := strconv.Atoi(vars["height"])
```

## 6.7 Error Handling (05:00)

- new error

```go
var ErrNotFound = errors.New("block not found")
```

- new errorResponse

```go
type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}
```

## 6.8 Middlewares (10:01)

- Middleware is a function to call before final destination

### adapter pattern

- `Handler` is an interface implementing method called ServerHTTP
- `HandlerFunc` is type (adapter)
  - `HandlerFunc()`: constructing a type
  - adaptor ask us to send correct argument and adaptor implement everything we need

# 7. CLI

## 7.0 Introduction

- flag
- cobra

## 7.1 Parsing Commands (05:52)

- os.Args gives array of commands

```sh
go run main.go someCMD
-> [.../exe/main someCMD]
```

- to exit, use `os.Exit(0)`

## 7.2 FlagSet (10:26)

- flagSet is useful if one command has many flags

```go
rest := flag.NewFlagSet("rest", flag.ExitOnError)
portFlag := rest.Int("port", 4000, "Sets the port of the server")
...
rest.Parse(os.Arge[2:])
...
if rest.Parsed() {
		fmt.Println(*portFlag)
		fmt.Println("Start server")
}
```

## 7.3 Flag (10:08)

- easier than flagSet

```go
port := flag.Int("port", 4000, "Set port of the server")
mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'"
```

- refactor main.go > cli/cli.go

```sh
mkdir cli
sed 's/main()/Start()/; s/package main/package cli/' main.go > cli/cli.go
echo "package main\nfunc main() {}" > main.go
```

- challenge: make command to run both with differen port and goroutine

# 8. PERSISTENCE

## 8.0 Introduction to Bolt (08:09)

- currently everthing is on memory (slice of block)
- bolt: key/value database specified for get/set
  - eg) "sdkfljsdlfjds": {"data: PrvHash"}

## 8.1 Creating the Database (11:47)

- There will be no immediate `Start` so that start coding from `db/db.go`

```sh
mkdir db
touch db/db.go
go get github.com/boltdb/bolt
```

```go
const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}
```

## 8.2 A New Blockchain (11:53)

- divide & conquer

```sh
mv blockchain/blockchain.go blockchain/chain.go
touch blockchain/block.go
```

## 8.3 Saving A Block (12:25)

- gob: encode/decode data<->byte
- buffer: place to put byte with write/read

```go
// blockchain/block.go
func (b *Block) toBytes() []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	utils.HandleErr(encoder.Encode(b))
	return blockBuffer.Bytes()
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, b.toBytes())
}
```

- SaveBlock: bolt can save only byte

```go
// db/db.go
func SaveBlock(hash string, data []byte) {
	fmt.Printf("Saving Block %s\nData: %b\n", hash, data)
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}
```

## 8.4 Persisting The Blockchain (10:55)

- move ToBytes from `block.go` to `utils.go`
  - inferface -> can get any type

```go
// utils/utils.go
func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}
```

## 8.5 Restoring the Blockchain (12:28)

- when we start, should restore chain from checkpoint
- restore from byte to data

```go
func (b *blockchain) restore(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	utils.HandleErr(decoder.Decode(b)) // with pointer, modify the origin value from byte to data
}
```

- select checkpoint

```go
func Checkpoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket(([]byte(dataBucket)))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}
```

## 8.6 Restoring Block (13:27)

- Add func FromBytes at `utils.go`
- Refactor func restore with func FromBytes at `chain.go`
- Add func Block at `db/db.go`
- Add ErrNotFound, func restore, func FindBlock at `block.go`
- Omit GET, POST case at `rest.go` for test

## 8.7 All Blocks (10:51)

- Close at `db/db.go`
- defer Close() at `main.go`
  - if there is Goexit, deferred calls will be executed
- Add func Blocks at `chain.go`
- Recover GET, POST case at `rest.go`

## 8.8 Recap (10:46)

- Refactor SaveBlockchain -> SaveCheckpoint
- bolt, get/set data to bucket
- singletone -> if no checkpoint -> create genesis
- func persist -> save to bolt database

# 9 MINING

## 9.0 Introduction to PoW (06:28)

- Proof Of Work
- Add properties to Block struct at `block.go`

```go
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
```

- Delete db

```sh
rm blockchain.db
```

## 9.1 PoW Proof Of Concept (10:13)

- Mining prototype

```go
difficulty := 2
target := strings.Repeat("0", difficulty)
nonce := 1
for {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte("hello"+fmt.Sprint(nonce))))
	fmt.Printf("Hash:%s\nTarget:%s\nNonce:%d\n\n", hash, target, nonce)
	if strings.HasPrefix(hash, target) {
		return
	} else {
		nonce++
	}
}
```

## 9.2 Mine Block (07:48)

- func mine at `block.go`
- Modify way get hash -> `block.mine()` at `block.go`

## 9.3 Difficulty part One (10:20)

- Refactor Hash from `block.go` to `utils.go`

```go
s := fmt.Sprintf("%v", i)
// v means default formatter
```

- Add timestamp on each mining at `block.go`

```go
b.Timestamp = int(time.Now().Unix())
```

- Add func difficulty at `chain.go`

```go
func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
	} else {
		return b.CurrentDifficulty
	}
}
```

## 9.4 Difficulty part Two (11:56)

- Add const blockInterval, allowedRange at `chain.go`
- Whenever create new block, `b.CurrentDifficulty = block.Difficulty`
- Adjust allowedRange rather than constant value

```go
if actualTime <= (expectedTime - allowedRange) {
	return b.CurrentDifficulty + 1
} else if actualTime >= (expectedTime + allowedRange) {
	return b.CurrentDifficulty - 1
}
return b.CurrentDifficulty
```

## 9.5 Conclusions (06:28)

- Add status route at `rest.go`
- Refactor Handling errors for encoder at `rest.go`
- Cheack status of difficulty grows each 5 blocks

# 10 TRANSACTIONS

## 10.0 Introduction (04:37)

- Course ~#9 were all about Protecting Data
- Here going to learn Moving value between our user
- uTxOut: Unspent Transaction Output mpdel?

## 10.1 Introduction to Transactions (07:01)

- Tx
  - TxIn[$5(me)]: money that i have
  - TxOut[$0(me), $5(you)]: money that every body has by the end of Tx
  - Just change owner
- What if currency scale is different? change?
  - TxIn[$10(me)]
  - TxOut[$5(me), $5(you)]
- Just find last Tx
- Coinbase input: created by blockchain, to miner
  - TxIn[$10(blockchain)]
  - TxOut[$10(miner)]

## 10.2 Coinbase Transaction (11:18)

```sh
touch blockchain/transactions.go
```

- Add struct Tx, TxIn, TxOut, func makeCoinbaseTx at `transactions.go`
- Remove Data, Add Transaction to struct Block at `block.go`
- Refactor to remove Data at `rest.go, explorer.go, block.go, chain.go`

## 10.3 Balances (13:30)

- Add func txOuts, TxOutsByAddress, BalanceByAddress at `chain.go`
- Add struct balanceResponse, func balance, router balance at `rest.go`

## 10.4 Mempool (04:28)

- Memory pool: where we put unconfirmed transaction
- After confirmed -> Part of Block

## 10.5 AddTx (10:42)

- Cuz Mempool is on the memory, we don't need to initialize
- Add struct mempool, var Mempool, func makeTx, AddTx at `transactions.go`
- Add router mempool at `rest.go`

## 10.6 makeTx (11:33)

- Improve func maskTx at `transactions.go`:
  - BalanceByAddress < amount -> error
  - Until total >= mount, append &TxIn{txOut.Owner, txOut.Amount} to txIns
  - change := total - amoun -> append to txOuts
  - append &TxOut{to, amount} to txOuts
- Add struct addTxPayload, router transactions at `rest.go`

## #10.7 Confirm Transactions (11:54)

- Refactor Transactions to after block.mine() at `block.go`

```go
block.mine()
block.Transactions = Mempool.TxToConfirm()
```

- Add func TxToConfirm at `transactions.go`
- Duplication bug: we should check if the coin of txOut was already used or not

## 10.8 uTxOuts (12:02)

- We need to find which one is duplicated

```
Tx1
	TxIns[COINBASE]
	TxOuts[$5(you)] <- Spent TxOut
Tx2
	TxIns[Tx1.TxOuts[0]]
	TxOuts[$5(me)] <- unspent TxOut to Spent
Tx3
	TxIns[Tx2.TxOuts[0]]
	TxOuts[$3(you), $2(me)] <- uTxOut * 2
```

- Modify and Add struct TxIn, UTxOut at `transactions.go` with

```go
TxID   string
Index  int
Amount int
```

- Remove func txOuts, Rename TxOutsByAddress to UTxOutsByAddress at `chain.go`

## 10.9 UTxOutsByAddress (10:37)

- Implement UTxOutsByAddress at `chain.go`
  - In all inputs, if address equal to the owner, append TxId to creatorTxs
  - In all outputs, if address is equal to the owner and TxId is `not` in creatorTxs, append to uTxOuts
- Refactor Tx.Id, getId => `ID` at `transactions.go`

## 10.10 makeTx part Two (10:04)

- Impl makeTx at `transactions.go`
- But, It stil copy the coin if the Tx is on mempool
  - it checkes spent or unspent only (confirmed)
  - Should check unconfirmed Tx too.

## 10.11 isOnMempool (06:55)

- Impl func isOnMempool at `transactions.go`
- Add !isOnMempool condition to UTxOutsByAddress at `chain.go`
- looks like working but why error message is "not enough funds" not "not enough money"?
  - does rest.go convert error automatically?

## 10.12 Refactoring (16:46)

### Refactor redundant loop

1. return true to kill func way

```go
func isOnMempool(uTxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				return true
			}
		}
	}
	return false
}
```

2. break Outer labeled loop way

```go
func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}
```

### Method Vs Function

- If it is mutating struct -> `method`
- Else, -> normal `func` with struct as input param
- Sort method first, func last

## 10.13 Deadlock (03:53)

- Current func Blockchain at `chain.go` is recursive
- Because no call to Do returns until the one call to f returns, if f causes Do to be called, it will deadlock.
- Modify logic
  - Remove nil condition of func Blockchain at `chain.go`
  - func createBlock receives diff param at `block.go`
  - Rename difficulty -> getDifficulty at `chain.go`

## 11.0 Introduction (04:50)

- If he owns unspent output
- If he approved the transaction

1. How signature, verification works
2. Persistance to db
3. Impl signature, verification with tracsaction

```sh
mkdir wallet
touch wallet/wallet.go
```

## 11.1 Private and Public Keys (08:26)

1. Hash the msg

```go
"i love you" -> hash(x) -> "hashed_message"
```

2. Generate key pair

```go
keypair (privateKey, publicKey)
	(save privateK to a file -> wallet)
```

3. Sign the hash

```go
("hashed_message" + privateKey) -> "signature"
```

4. Verify with publicKey

```go
("hashed_message" + "signature" + publicKey) -> true / false
```

## 11.2 Signing Messages (10:25)

- ecdsa: Elliptic Curve Digital Signature Algorithm

```go
privateKey = ecdsa.GenerateKey
hashAsBytes = hex.DecodeString(Hash(message))
r, s = ecdsa.Sign
```

## 11.3 Verifying Messages (13:23)

- Refactor generated values to constant var

```go
privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // bigint
keyAsBytes, err := x509.MarshalECPrivateKey(privateKey) // bigint -> hex
utils.HandleErr(err)
fmt.Printf("privateKey = %x\n\n", keyAsBytes)

message := "i love you"
hashedMessage := utils.Hash(message)
fmt.Printf("hashedMessage = %s\n\n", hashedMessage)

hashAsBytes, err := hex.DecodeString(hashedMessage)
utils.HandleErr(err)
r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
signature := append(r.Bytes(), s.Bytes()...) // [32]byte + [32]byte
fmt.Printf("signature = %x\n\n", signature) // [64]byte -> hex

utils.HandleErr(err)
ok := ecdsa.Verify(&privateKey.PublicKey, hashAsBytes, r, s)
fmt.Println(ok)
```

```go
const (
	privateKey    string = "30770201010420270623da3768df6fc3c3439b8e0319621318b1dec6199052f49faefdd9d80548a00a06082a8648ce3d030107a1440342000462ded99b11da850eec19a908aa57effbec88541aa04da07d0a2cabf046b2502dd061eccc9860c7922ea758a2e8ac1e5f6d044d7a6af03060aa5dcb13cafc8a73"
	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
	signature     string = "6d56582490ff9a54b44df6bf9fa991c0432f2fd25f32760bf540b10049b50a048ca986d7e9f0ee7745bce735dcd0db951f21664f054e94b0a03d87046857a3ca"
)
```

## 11.4 Restoring (12:26)

```go
privByte, err := hex.DecodeString(privateKey)
utils.HandleErr(err)

privateKey, err := x509.ParseECPrivateKey(privByte)
utils.HandleErr(err)

sigBytes, err := hex.DecodeString(signature)

rBytes := sigBytes[:len(sigBytes)/2]
sBytes := sigBytes[len(sigBytes)/2:]

var bigR, bigS = big.Int{}, big.Int{}
bigR.SetBytes(rBytes)
bigS.SetBytes(sBytes)

hashBytes, err := hex.DecodeString(hashedMessage)
utils.HandleErr(err)
ok := ecdsa.Verify(&privateKey.PublicKey, hashBytes, &bigR, &bigS)
fmt.Println(ok)
```

## 11.5 Wallet Backend (10:41)

- Add func hasWalletFile at `wallet.go`
- Singletone pattern func Wallet skeleton at `wallet.go`

## 11.6 Persit Wallet (09:16)

- Add func createPrivKey, persistKey at `wallet.go`
- If there is new var togather, we can recreate err with `newVar err :=`
  - (actually updating)

1. pure update way

```go
bytes, err := x509.MarshalECPrivateKey(key)
utils.HandleErr(err)
err = os.WriteFile(fileName, bytes, 0644)
...
```

2. recreate way

```go
bytes, err := x509.MarshalECPrivateKey(key)
utils.HandleErr(err)
newVar2, err := os.WriteFile(fileName, bytes, 0644)
```

## 11.7 Restore Wallet (09:15)

- Add func restoreKey at `wallet.go`
- Names return: good for short func, easy to understand only with signature, bad for long function's retunry `empty` and lookup above again

```go
func restoreKey() (key *ecdsa.PrivateKey ){
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}
```

## 11.8 Addresses (09:20)

- Add func aFromK, sign at `wallet.go`
- Replace `"nico"` to `wallet.Wallet().Address` at `transactions.go`

## 11.9 Verification Function (11:24)

- Add func restoreBigInts, verify at `wallet.go`
- restoreBigints can handle signature, address(to publicKey)

## 11.10 Recap (08:56)

- Add func encodeBigInts at `wallet.go`
- Refactor with encodeBigInts

## 11.11 Transaction Signing (11:01)

- Add func Txs, FindTx at `chain.go`
- Rename TxIn.Owner -> Signature, TxOut.Owner -> Address at `transactions.go`
- Add func sign at `transactions.go`

## 11.12 Transaction Verification (16:36)

### `transactions.go`

#### Impl func `validate`

- To validate New Tx(in),
  - Find prev Txout referencing same TxID with New Tx
  - Check if address is same (same owner)
  - Verify signature created by private key of the owner

#### Modify func `makeT`

- check valid, return ErrorNotValid

### `chain.go`

#### Modify func `UTxOutsByAddress`

```go
input.Onwer == address
```

⬇️

```go
if FindTx(b, input.TxID).TxOuts[input.Index].Address == address
```

## 11.13 Conclusions (06:58)

### `rest.go`

#### Add router `wallet`

#### Add func myWallet with Anonymous struct

```go
func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(struct {
		Address string `json:"address"`
	}{Address: address})

}
```

### Q

- Don't we need to use address when we send coin instead of name?
- Study more (signature - address - name) relationship
- How to filter mempool-ed Tx?

## 11.14 Recap (10:16)

### `rest.go`

#### func transactions: Superfluous response.WriteHeader error

- Should return http.StatusBadRequest(404)
- Should use err.Error() to convey original error instead of literal error

```go
json.NewEncoder(rw).Encode(errorResponse{err.Error()})
```

- return to finish function

### Recap

```
(TxOut1(publickKey), TxOut2)

Tx
	TxIn[
		(TxOut1)
		(TxOut2)
	]
	Sign with my privateKey

TxIn.Sign + TxOut1.Address -> true / false
```

## 12.0 Introduction (05:55)

- Lean Peer To Peer by simple Chatting app

## 12.1 Why Go Routines (09:44)

### Goroutine

- Running function at separate parallel dimension
- Can't assign or return value to variable immediately;

```go
go countToTen()
```

## 12.2 Channels (16:07)

### Channels

- Deadlock: Channel should not receive more than coroutine

## 12.3 Raead, Receive and Close (11:00)

```go
func countToTen(c chan<- int) { // send only chan<-
	for i := range [10]int{} {
		time.Sleep(1 * time.Second)
		fmt.Printf("sending %d\n", i)
		c <- i
	}
	close(c)
}

func receive(c <-chan int) { // receive only <-chan
	for {
		a, ok := <-c // blocking by getting next value
		if !ok {
			fmt.Println("Done")
			break
		}
		fmt.Printf("received %d\n", a)
	}
}

func main() {
	c := make(chan int)
	go countToTen(c)
	receive(c)
}
```

## 12.4 Buffered Channels (14:29)

```go
c := make(chan int, NumOfBuffer)
```

- Don't block first N values, then block/wait the queue like normal Unbuffered channel

## 12.5 WebSocket Upgrades (11:47)

- HTTP: stateless
- WS: statefull, connected
- Upgrade Go with WS

```sh
mkdir p2p
touch p2p/p2p.go
touch chat.html
go get github.com/gorilla/websocket
go run -race main.go run -mode=rest -port=3000
```

### `chat.html`

```html
const socket = new WebSocket("ws://localhost:4000/ws");
```

### `p2p/p2p.go`

- Add func Upgrade

### `rest/rest.go`

- Add ws to URL, router
- Add func loggerMiddleware

## 12.6 ReadMessage (13:18)

### `p2p/p2p.go`

#### func `Upgrade`

- Add CheckOrigin
- Add conn.ReadMessage with for loop
  - conn.ReadMessage is a receiver Channel for socket

## 12.7 WriteMessage (11:02)

### `chat.html`

- Add form, send, receive eventListener

### `p2p/p2p.go`

- Add conn.WriteMessage

## 12.8 Connections (13:02)

- Connect client to client through server
- http.ListenAndServe uses goroutine

### `p2p/p2p.go`

- Add slice of conn, append conn, send to other conns,

### To-do

- If browser is refreshed, gets error -> How to handle closed connection?
- A Message should not block others -> How to separate?

## 12.9 Peers (17:51)

- Should connect peer and peer not through server

### `p2p/p2p.go`

- Add func `AddPeer`

### `rest.go`

- Add struct addPeerPayload, func peers, router peers

## 12.10 initPeer (10:29)

### `p2p/peer.go`

```sh
touch p2p/peer.go
```

- Add struct `peer`, func `initPeer`

### `rest.go`

- Add GET to router peers

## 12.11 openPort (10:46)

### `p2p/p2p.go`

#### func AddPeer

- send openPort with URL

#### func Upgrade

- get openPort from URL

### `rest.go`

#### func peers

- p2p.Addpeer(..., port)

## 12.12 Recap (14:24)

### `utils.go`

- Add func Spliiter

### `p2p.go`

- Refactor with utils.Splitter

## #12.13 Read Peer (11:06)

### `p2p/peer.go`

- Add method read
- Add go p.read() at func `initPeer`

### `p2p/p2p.go`

- Add conn.WriteMessage at func `Upgrade`

## 12.14 Inbox (08:30)

- Instead of writing once, use coroutine and channel

### `p2p/peer.go`

- Add method `write`

#### func `initPeer`

- make channel initiating peer
- go p.write()

### `p2p/p2p.go`

- send message to inbox(channel)

## 12.15 Cleanup (09:57)

### `p2p/p2p.go`

- Delete Hellos, Add `initPeer()`

### `p2p/peer.go`

- Add key, address, port to struct peer
- Add func `close`
- Add defer p.close() at `read`, `write`

## 12.16 Data Race (12:39)

- Stable bolt has socket hangs up error with race
- Change to bbolt

```sh
go get go.etcd.io/bbolt
```

```go
// db.go
import (bolt "go.etcd.io/bbolt")
```

- Data Race: When more than two goroutine access to the same block
- When we read and modify peers at the same time, gets Data race error

## 12.17 Mutex (13:44)

### `p2p/peer.go`

- Add type peers with `sync.Mutex`

```go
type peers struct {
	v map[string]*peer
	m sync.Mutex
}

var Peers peers = peers{
	v: make(map[string]*peer),
}
```

- Should lock before read of delete

```go
func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.v, p.key)
}
```

- Add func `AllPeers` and convert output from object to array of keys
  - Cuz we don't modify peers but just read -> not method but func -> use getter

```json
{
  "127.0.0.1:4000": {},
  "...": {}
}
```

⬇️

```json
[
	"127.0.0.1:4000",
	...
]
```

## 12.18 Demonstration (08:44)

- Can use func with defer to delay between lock and unlock
  - -> Can Demonstrate the mutex

```go
func (p *peer) close() {
	Peers.m.Lock()
	defer func() {
		time.Sleep((20 * time.Second))
		Peers.m.Unlock()
	}()
	p.conn.Close()
	delete(Peers.v, p.key)
}
```

## 12.19 Messages (13:41)

### Messaging Scenario

#### Behind case

1. :4000 sends newest block to :3000
2. :3000 realizes that :3000 is behind, ask to :4000 latest all blocks
3. :4000 gives latest all blocks

#### Ahead case

1. :4000 sends newest block to :3000
2. :3000 realizes that :3000 is ahead, send newest block to :4000
3. :4000 realizes that :4000 is behind, ask to :3000 latest all blocks
4. :3000 gives latest all blocks

### `p2p/messages.go`

```sh
touch p2p/messages.go
```

### `iota`

- Auto increament sequence

```go
const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)
```

### `p2p/peer.go`

- Let's Communicate with json
- Change from `ReadMessage` to `ReadJSON`

## 12.20 Newest Block (11:31)

### `p2p/messages.go`

- Add func addPayload, makeMessage, sendNewestBlock
- Why should we json.Marshal twice?
  - cuz Payload is public, it can be already JSON
  - Both Payload and Message should be json.Marshal-ed

### `p2p/p2p.go`

- sendNewstBlock in addPeer

## 12.22 Handle Message (09:50)

### `utils.go`

- Add func ToJSON

### `p2p/messages.go`

- Add func handleMsg with Unmarshal
- Remove func addPayload, Use ToJSON instead.

### `p2p/peer.go`

- handleMsg whenever read()

## 12.23 Recap (10:01)

### `p2p/peer.go`

- Datarace: Though there is lock at close(), opposite side stil add Peer twice -> Add Lock at initPeer

## 12.24 All Blocks (14:19)

### `db/db.go`

- Each client should use different db
- Add func getDbName: get db name with port number

### `p2p/messages.go`

- Add func requestAllBlocks, sendAllBlocks

#### func `handleMsg`

- if `payload.Height >= b.Height` requestAllBlocks()
  - To handle case when height is equal, used some illgical trick
- Add case `MessageAllBlocksRequest`, `MessageAllBlocksresponse`

- Why don't we just get height with blockchain.Height?

## 12.25 Recap (10:15)

- Add console print

### Case1

1. 4000: wants to connect to port 3000
2. 3000: 4000 wants an upgrade
3. 4000: Sending newest block to 127.0.0.1:3000
4. 3000: Received the newest block from 127.0.0.1:4000
5. 4000: 127.0.0.1:3000 wants all the blocks.
6. 3000: Received all the blocks from 127.0.0.1:4000

### Case2

- After mining blocks at 3000,

1. 4000: wants to connect to port 3000
2. 3000: 4000 wants an upgrade
3. 4000: Sending newest block to 127.0.0.1:3000
4. 3000: Received the newest block from 127.0.0.1:4000 // 3000 have more than 4000
5. 3000: Sending newest block to 127.0.0.1:4000
6. 4000: Received the newest block from 127.0.0.1:3000 // realize 4000 is behind 3000
7. 4000: Requesting all blocks to 127.0.0.1:3000
8. 3000: 127.0.0.1:4000 wants all the blocks.
9. 4000: Received all the blocks from 127.0.0.1:3000

## 12.26 Replace Blockchain (10:59)

- Can syncronize and persist now!

### `/blockchain/block.go`

- Rename method persist -> func persitBlock

### `/blockchain/chain.go`

- Add method Replace
  - Renew blockchain
  - EmptyBlocks and persist newBlocks

### `db/db.go`

- Add funcEmptyBlocks()

### `p2p/messages.go`

- blockchain.Blockchain().Replace(payload) at the end of handleMsg

## 12.27 Broadcast Block (11:29)

### Solve datarace: cover blockchain with mutex

#### `blockchain/chain.go`

- Add m to struct blockchain
- Lock & unlock at func Blocks, Status, Replace

### Broadcast: send same msg to everybody

#### `blockchain/chain.go`

- func AddBlock returns newBlock
- Add func Status: show status of blockchain at `/status`

#### `p2p/messages.go`

- Add iota MessageNewBlockNotify
- Add func notifyNewBlock, BroadcastNewBlock

#### `rest/rest.go`

- `/blocks`: after AddBlock, BroadcastNewBlock
- `/status`: blockchain.Status
