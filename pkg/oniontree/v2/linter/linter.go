package linter

import (
	"encoding/json"
	"errors"
	"github.com/onionltd/oniontree-tools/pkg/oniontree/v2"
	"github.com/onionltd/oniontree-tools/pkg/oniontree/v2/linter/schema"
	"github.com/xeipuuv/gojsonschema"
)

type Linter struct{}

func (l Linter) Lint(service oniontree.Service) error {
	b, err := json.Marshal(service)
	if err != nil {
		return err
	}

	schemaStr := schema.V0

	schemaLoader := gojsonschema.NewStringLoader(schemaStr)
	documentLoader := gojsonschema.NewBytesLoader(b)

	res, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !res.Valid() {
		errs := []error{}
		for _, errMsg := range res.Errors() {
			if errMsg == nil {
				continue
			}
			errs = append(errs, errors.New(errMsg.String()))
		}
		return LinterError{errs}
	}

	return nil
}

func NewLinter() *Linter {
	return &Linter{}
}
