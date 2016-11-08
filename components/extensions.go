package components

import (
	"errors"
	"strings"
)

type errorList []error

func (errs errorList) asError() error {
	if (errs == nil) || len(errs) == 0 {
		return nil
	}

	var results []string

	for _, e := range errs {
		results = append(results, e.Error())
	}

	result := strings.Join(results, "\n")

	return errors.New(result)
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}

	return y
}
