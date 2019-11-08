package main

import (
	"flag"
	"fmt"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"io/ioutil"
	"os"
)

func exitError(msg string) {
	fmt.Printf("%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}

func main() {
	id := flag.String("id", "", "Onion service ID.")
	description := flag.String("description", "", "Public key description.")
	file := flag.String("file", "", "Public key file.")
	replace := flag.Bool("replace", false, "Replace existing public keys.")
	flag.Parse()

	if *id == "" {
		exitError("id not specified")
	}
	if *file == "" {
		exitError("file not specified")
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

	sk, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	k, err := service.ParseKey(sk)
	if err != nil {
		panic(err)
	}
	k.Description = *description

	if *replace {
		s.SetPublicKeys(k)
	} else {
		s.AddPublicKeys(k)
	}

	if err := onionTree.Edit(*id, s); err != nil {
		panic(err)
	}
}
