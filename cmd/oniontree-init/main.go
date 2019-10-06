package main

import (
	"flag"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"os"
)

func main() {
	flag.Parse()

	dir := flag.Arg(0)
	if dir == "" {
		var err error
		if dir, err = os.Getwd(); err != nil {
			panic(err)
		}
	}
	onionTree := oniontree.New(dir)
	if err := onionTree.Init(); err != nil {
		panic(err)
	}
}
