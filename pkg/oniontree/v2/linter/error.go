package linter

import "strings"

type LinterError struct {
	errs []error
}

func (e LinterError) Error() string {
	s := make([]string, 0, len(e.errs))
	for i := range e.errs {
		s = append(s, e.errs[i].Error())
	}
	return strings.Join(s, "\n")
}
