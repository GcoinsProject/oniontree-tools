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
	name := flag.String("name", "", "Onion service name.")
	description := flag.String("description", "", "Onion service description.")
	urls := flag.String("urls", "", "Onion service URL.")
	replace := flag.Bool("replace", false, "Replace URLs instead of adding them.")
	flag.Parse()

	if *id == "" {
		exitError("id not specified")
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	onionTree, err := oniontree.Open(wd)
	if err != nil {
		panic(err)
	}

	s, err := onionTree.Get(*id)
	if err != nil {
		if err == oniontree.ErrIdNotExists {
			exitError(err.Error())
		}
		panic(err)
	}

	if *name != "" {
		s.Name = *name
	}
	if *description != "" {
		s.Description = *description
	}
	if *urls != "" {
		if *replace {
			s.SetURLs(strings.Split(*urls, ",")...)
		} else {
			s.AddURLs(strings.Split(*urls, ",")...)
		}
	}

	if err := onionTree.Edit(*id, s); err != nil {
		panic(err)
	}
}
