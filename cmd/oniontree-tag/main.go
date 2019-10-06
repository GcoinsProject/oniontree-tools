package main

import (
	"flag"
	"fmt"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"os"
)

func exitError(msg string) {
	fmt.Printf("%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}

func main() {
	id := flag.String("id", "", "Onion service ID.")
	tag := flag.String("tag", "", "Tag name.")
	flag.Parse()

	if *id == "" {
		exitError("id not specified")
	}
	if *tag == "" {
		exitError("tag not specified")
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	onionTree := oniontree.New(wd)
	if err := onionTree.Tag(*id, []string{*tag}); err != nil {
		if err == oniontree.ErrIdNotExists {
			exitError(err.Error())
		}
		panic(err)
	}
}
