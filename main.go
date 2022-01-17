package main

import (
	"fmt"

	"github.com/devgony/nomadcoin/person"
)

func main() {
	henry := person.Person{}
	henry.SetDetails("henry", 12)
	fmt.Println(henry)
}
