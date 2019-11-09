package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"github.com/xeipuuv/gojsonschema"
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

//go:generate go run generate.go

func main() {
	id := flag.String("id", "", "Onion service ID.")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	onionTree, err := oniontree.Open(wd)
	if err != nil {
		panic(err)
	}

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
		b, err := onionTree.GetRaw(id)
		if err != nil {
			panic(err)
		}

		// Validate YAML syntactically
		s := &service.Service{}
		if err := yaml.Unmarshal(b, &s); err != nil {
			exitnum = 1
			printLintStatus(id, err)
			continue
		}

		// Convert service data to JSON so it can be validated by jsonschema
		b, err = json.Marshal(s)
		if err != nil {
			panic(err)
		}

		schemaLoader := gojsonschema.NewStringLoader(ServiceSchema)
		documentLoader := gojsonschema.NewBytesLoader(b)
		res, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			panic(err)
		}

		if !res.Valid() {
			for _, errMsg := range res.Errors() {
				printLintStatus(id, errors.New(errMsg.String()))
			}
			continue
		}
	}
	os.Exit(exitnum)
}
