package main

import (
	"flag"
	"fmt"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
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
	tags := flag.String("tags", "", "Onion service tags.")
	flag.Parse()

	if *id == "" {
		exitError("id not specified")
	}
	if *urls == "" {
		exitError("url not specified")
	}

	s := service.Service{}
	s.Name = *name
	s.Description = *description
	s.SetURLs(strings.Split(*urls, ",")...)

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	onionTree, err := oniontree.Open(wd)
	if err != nil {
		panic(err)
	}

	if err := onionTree.Add(*id, s); err != nil {
		if err == oniontree.ErrIdExists {
			exitError(err.Error())
		}
		panic(err)
	}
	if *tags != "" {
		if err := onionTree.Tag(*id, strings.Split(*tags, ",")); err != nil {
			panic(err)
		}
	}
}
