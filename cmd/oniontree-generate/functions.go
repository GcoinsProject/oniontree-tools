package main

import (
	"encoding/json"
	"text/template"
)

var functions = template.FuncMap{
	// Convert an object to JSON.
	"json": func(v interface{}) string {
		b, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(b)
	},
	// Convert an object to JSON with indentation.
	"jsonPretty": func(v interface{}) string {
		b, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			panic(err)
		}
		return string(b)
	},
}
