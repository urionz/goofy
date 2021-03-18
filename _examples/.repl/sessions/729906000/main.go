package main

import (
	"encoding/json"
	"fmt"
)

var (
	bs, _ = json.Marshal("1")
)

func main() {
	fmt.Printf("<%T> %+v\n", echo, echo)
}
