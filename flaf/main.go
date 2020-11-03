package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("=============")
	for _, v := range os.Args {
		fmt.Println(v)
	}
	fmt.Println("=============")
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe != 0 {
		acl, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println(string(acl))
	}
	fmt.Println("=============")
}
