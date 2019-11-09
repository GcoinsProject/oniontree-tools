// +build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"
)

const fileTemplate = `// AUTO GENERATED FILE, DO NOT EDIT!
package main

const ServiceSchema = {{ printf "%q" . }}
`

func main() {
	t := template.Must(template.New("").Parse(fileTemplate))

	fIn, err := os.Open("schema.json")
	if err != nil {
		panic(err)
	}
	defer fIn.Close()

	b, err := ioutil.ReadAll(fIn)
	if err != nil {
		panic(err)
	}

	v := make(map[string]interface{})
	if err := json.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	b, err = json.Marshal(v)
	if err != nil {
		panic(err)
	}

	fOut, err := os.Create("schema.go")
	if err != nil {
		panic(err)
	}
	defer fOut.Close()

	if err := t.Execute(fOut, b); err != nil {
		panic(err)
	}
}
