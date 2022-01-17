# 2.2 Why Go (06:17)

- Only 1 option: for loop
- Powerful standard library
- Low learning curve

# 3.0 Creating the Project (04:30)

```sh
go mod init github.com/devgony/nomadcoin
ls
> go.mod # like package.json
touch main.go
```

# 3.1 Variables in Go (07:54)

## create and update var syntax (should be inside func only)

```go
var name string = "henry"
name := "henry" // same with above, syntax sugar
```

- var, const
- bool, string, int(8,16..64), uint(8,16..64) byte, float(32,64)

# 3.2 Functions (08:59)

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

# 3.3 fmt (03:52)

```go
x := 84375983402
fmt.Printf("%b\n", x)
fmt.Printf(fmt.Sprintf("%b\n", x)) // return fmted string (not print)
fmt.Printf("%o\n", x)
fmt.Printf("%x\n", x)
fmt.Printf("%U\n", x)
```

# 3.4 Slices and Arrays (08:02)

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

# 3.5 Pointers (08:44)

```go
a := 2
b := a  // copy
c := &a // borrow
a = 9
fmt.Println(a, b, *c) // 9 2 9
```

# 3.6 Structs (07:11)

## struct (like class)

```go
type person struct {
	name string
	age  int
}
```

## receiver function (like methods)

- choose first letter of type as a alias

```go
func (p person) sayHello() {
	// give access to instance
	fmt.Printf("Hello! My name is %s and I'm %d", p.name, p.age)
}
```

# 3.7 Structs with Pointers (09:35)

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

## Receiver pointer function

- basically receiver function's p is copy -> use only for read
- we don't want to mutate copied, but origin
- use \* to mutate original

```go
func (p *Person) SetDetails(name string, age int) {
    // p is the origin
	p.name = name
	p.age = age
```

# 4.0 Introduction (05:05)

- Concentrate to blockchain concept, solve side problem later

# 4.1 Our First Block (13:58)

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

# 4.2 Our First Blockchain (09:43)

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

# 4.3 Singleton Pattern (05:57)

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
