package main

import (
	"flag"
	"fmt"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"os"
	"strings"
)

func exitError(msg string) {
	fmt.Printf("%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}

func main() {
	id := flag.String("id", "", "Onion service ID.")
	tags := flag.String("tags", "", "Onion service tags.")
	flag.Parse()

	if *id == "" {
		exitError("id not specified")
	}
	if *tags == "" {
		exitError("tag not specified")
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	onionTree := oniontree.New(wd)

	if err := onionTree.Tag(*id, strings.Split(*tags, ",")); err != nil {
		if err == oniontree.ErrIdNotExists {
			exitError(err.Error())
		}
		panic(err)
	}
}
