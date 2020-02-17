package main

import (
	"flag"
	"fmt"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"os"
	"text/template"
)

func exitError(msg string) {
	fmt.Printf("%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}

type listerConfig struct {
	IncludeDescription bool
	IncludePublicKeys  bool
}

func main() {
	cfg := listerConfig{}
	onionTreePath := flag.String("oniontree", "", "Path to OnionTree directory")
	templatePath := flag.String("template", "", "Path to template")
	flag.BoolVar(&cfg.IncludeDescription, "with-description", true, "Include service description")
	flag.BoolVar(&cfg.IncludePublicKeys, "with-public-keys", true, "Include public keys")
	flag.Parse()

	if *onionTreePath == "" {
		exitError("path to OnionTree not specified")
	}
	if *templatePath == "" {
		exitError("path to template not specified")
	}

	onionTree, err := oniontree.Open(*onionTreePath)
	if err != nil {
		panic(err)
	}

	templateData := struct {
		Unsorted  map[string]service.Service `json:"unsorted"`
		Tagged    map[string][]string        `json:"tagged"`
		Addresses map[string]string          `json:"addresses"`
	}{
		listUnsorted(onionTree, &cfg),
		listTagged(onionTree),
		listAddresses(onionTree),
	}

	if err := template.Must(template.New(*templatePath).
		Funcs(functions).
		ParseFiles(*templatePath)).
		Execute(os.Stdout, templateData); err != nil {
		panic(err)
	}
}

func listUnsorted(onionTree *oniontree.OnionTree, cfg *listerConfig) map[string]service.Service {
	serviceIDs, err := onionTree.List()
	if err != nil {
		panic(err)
	}
	resultMap := make(map[string]service.Service)
	for idx := range serviceIDs {
		s, err := onionTree.Get(serviceIDs[idx])
		if err != nil {
			panic(err)
		}
		modifyServiceData(&s, cfg)
		resultMap[serviceIDs[idx]] = s
	}
	return resultMap
}

func listTagged(onionTree *oniontree.OnionTree) map[string][]string {
	tagIDs, err := onionTree.ListTags()
	if err != nil {
		panic(err)
	}
	resultMap := make(map[string][]string)
	for _, tagID := range tagIDs {
		t, err := onionTree.GetTag(tagID)
		if err != nil {
			panic(err)
		}
		for _, serviceID := range t.Services {
			resultMap[tagID] = append(resultMap[tagID], serviceID)
		}
	}
	return resultMap
}

func listAddresses(onionTree *oniontree.OnionTree) map[string]string {
	serviceIDs, err := onionTree.List()
	if err != nil {
		panic(err)
	}
	resultMap := make(map[string]string)
	for idx := range serviceIDs {
		s, err := onionTree.Get(serviceIDs[idx])
		if err != nil {
			panic(err)
		}
		for _, url := range s.URLs {
			resultMap[url] = serviceIDs[idx]
		}
	}
	return resultMap
}

func modifyServiceData(s *service.Service, cfg *listerConfig) {
	if !cfg.IncludePublicKeys {
		s.PublicKeys = []service.PublicKey{}
	}
	if !cfg.IncludeDescription {
		s.Description = ""
	}
}
