package components

import (
	"errors"
	"strings"
)

type ErrorList []error

func (errs ErrorList) AsError() error {
	if len(errs) == 0 {
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
