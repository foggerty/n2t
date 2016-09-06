package components

import (
	"errors"
	"fmt"
)

type SymbolTable struct {
	variables map[string]int
	labels    map[string]int
}

func (st *SymbolTable) AddLabel(sym string, mem int) error {
	if _, ok := st.labels[sym]; ok {
		msg := fmt.Sprintf("Internal error - %s has already been added to the symbol table.", sym)
		return errors.New(msg)
	}

	st.labels[sym] = mem

	return nil
}

func (st *SymbolTable) AddVariable(sym string, mem int) error {
	if _, ok := st.variables[sym]; ok {
		msg := fmt.Sprintf("Internal error - %s has already been added to the symbol table.", sym)
		return errors.New(msg)
	}

	st.variables[sym] = mem

	return nil
}

func (st *SymbolTable) LabelLocation(sym string) (int, error) {
	if res, ok := st.labels[sym]; ok {
		return res, nil
	}

	msg := fmt.Sprintf("Internal error - %s not found in the symbol table.", sym)
	return 0, errors.New(msg)
}

func (st *SymbolTable) VariableLocation(sym string) (int, error) {
	if res, ok := st.variables[sym]; ok {
		return res, nil
	}

	msg := fmt.Sprintf("Internal error - %s not found in the symbol table.", sym)
	return 0, errors.New(msg)
}
