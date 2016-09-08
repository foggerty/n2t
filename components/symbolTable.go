package components

import (
	"errors"
	"fmt"
)

/*

 Note to self: everything here is passed by value.  But remember, that
 a map is a reference type, in that it's a small strut, part of which
 is a pointer back to a much larger data structure.  i.e. while the
 struct is always copied by value, it's values happen to be two
 automatically dereferenced pointers.

 Since SymbolTable is just two reference types (i.e. small struct with
 a pointer somewhere), there's not much point in making all of these
 methods work again pointers.  This includes the 'addToMap' and
 'getValue' functions.

 Also note that if AddLabel was to be written like so:

   func (st *SymbolTable) AddLabel(...

 everything still works as is, because when using the '.' operator
 against a pointer to a struct, the pointer is automatically
 dereferenced for you.  i.e, no need for a second '->' operator as in
 C.

*/

type SymbolTable struct {
	variables map[string]int
	labels    map[string]int
}

func NewSymbolTable() SymbolTable {
	return SymbolTable{
		variables: make(map[string]int),
		labels:    make(map[string]int)}
}

func (st SymbolTable) AddLabel(sym string, mem int) error {
	return addToMap(st.labels, sym, mem)
}

func (st SymbolTable) AddVariable(sym string, mem int) error {
	return addToMap(st.variables, sym, mem)
}

func addToMap(m map[string]int, s string, i int) error {
	if _, ok := m[s]; ok {
		msg := fmt.Sprintf("Internal error - %s has already been added to the table.", s)
		return errors.New(msg)
	}

	m[s] = i

	return nil
}

func (st SymbolTable) LabelLocation(sym string) (int, error) {
	return location(st.labels, sym)
}

func (st SymbolTable) VariableLocation(sym string) (int, error) {
	return location(st.variables, sym)
}

func location(m map[string]int, s string) (int, error) {
	if res, ok := m[s]; ok {
		return res, nil
	}

	msg := fmt.Sprintf("Internal error - %s not found in the lookup table.", s)

	return 0, errors.New(msg)
}
