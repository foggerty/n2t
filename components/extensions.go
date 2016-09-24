package components

import (
	"errors"
	"fmt"
	"strings"
)

type errorList []error

func (errs errorList) asError() error {
	if len(errs) == 0 {
		return nil
	}

	var results []string

	for e := range errs {
		results = append(results, fmt.Sprintf("%q\n", e))
	}

	result := strings.Join(results, "\n")

	return errors.New(result)
}
