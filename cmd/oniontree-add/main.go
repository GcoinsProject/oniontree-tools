package main

import (
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
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
	urls := flag.String("url", "", "Onion service URL")
	tags := flag.String("tags", "", "Onion service tags.")
	pubkey := flag.String("pubkey", "", "Admin's PGP key.")
	flag.Parse()

	s := service.Service{}
	s.Name = *name
	s.Description = *description
	s.PublicKey = *pubkey
	s.URLs = []string{}
	if *urls != "" {
		s.URLs = append(s.URLs, *urls)
	}

	if *id == "" {
		exitError("id not specified")
	}

	b, err := yaml.Marshal(s)
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	onionTree := oniontree.New(wd)

	if err := onionTree.Add(*id, b); err != nil {
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
