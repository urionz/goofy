package main

import "fmt"

type abc struct {
	bbb string
}

func somefun() string {
	return "abs"
}

var (
	m = &abc{bbb: "asdsdsa"}
)

func main() {
	fmt.Printf("<%T> %+v\n", m.bbb, m.bbb)
}
