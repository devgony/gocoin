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
