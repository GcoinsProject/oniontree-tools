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

func printLintStatus(file string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("%s: %s\n", file, err)
}

func main() {
	id := flag.String("id", "", "Onion service ID.")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	onionTree := oniontree.New(wd)

	ids := []string{}
	if *id == "" {
		var err error
		ids, err = onionTree.List()
		if err != nil {
			panic(err)
		}
	} else {
		ids = append(ids, *id)
	}

	exitnum := 0
	for _, id := range ids {
		b, err := onionTree.Get(id)
		if err != nil {
			panic(err)
		}

		s := &service.Service{}
		if err := yaml.Unmarshal(b, &s); err != nil {
			exitnum = 1
			printLintStatus(id, err)
			continue
		}
		// TODO: validate yaml schema!
	}
	os.Exit(exitnum)
}
