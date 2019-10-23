package main

import (
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"os"
)

func exitError(msg string) {
	fmt.Printf("%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}

func main() {
	id := flag.String("id", "", "Onion service ID.")
	name := flag.String("name", "", "Onion service name.")
	description := flag.String("description", "", "Onion service description.")
	urls := flag.String("url", "", "Onion service URL.")
	flag.Parse()

	sNew := service.Service{}
	sNew.Name = *name
	sNew.Description = *description
	if *urls != "" {
		sNew.AddURLs(*urls)
	}

	if *id == "" {
		exitError("id not specified")
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	onionTree := oniontree.New(wd)
	b, err := onionTree.Get(*id)
	if err != nil {
		if err == oniontree.ErrIdNotExists {
			exitError(err.Error())
		}
		panic(err)
	}

	sOld := service.Service{}
	if err := yaml.Unmarshal(b, &sOld); err != nil {
		panic(err)
	}

	if sNew.Name == "" {
		sNew.Name = sOld.Name
	}
	if sNew.Description == "" {
		sNew.Description = sOld.Description
	}
	if len(sNew.URLs) == 0 {
		sNew.URLs = sOld.URLs
	}
	sNew.PublicKeys = sOld.PublicKeys

	b, err = yaml.Marshal(sNew)
	if err != nil {
		panic(err)
	}

	if err := onionTree.Edit(*id, b); err != nil {
		panic(err)
	}
}
