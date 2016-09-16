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

type symbolTable struct {
	symbols     map[string]int
	initialised bool
}

func newSymbolTable() symbolTable {
	return symbolTable{
		symbols: make(map[string]int),
	}
}

func (st symbolTable) addLabel(sym string, mem int) {
	st.symbols[sym] = mem
}

func (st symbolTable) addVariable(sym string) {
	// To avoid confusion with labels, save with a default of -1.  This
	// is in case a label @LOOP reference comes before a label (LOOP)
	// statement.  Note that when a label IS saved, via addLabel(), the
	// initial value of -1 will be overwritten.
	st.symbols[sym] = -1
}

func (st symbolTable) symbolValue(sym string) (int, error) {
	if res, ok := st.symbols[sym]; ok {
		if res == -1 {
			panic("PROGRAMMER ERROR! - Symbol table was not initialised.")
		}

		return res, nil
	}

	msg := fmt.Sprintf("Internal error - %s not found in the lookup table.", sym)

	return 0, errors.New(msg)
}

func (st symbolTable) init() {
	mem := 16

	for k, v := range st.symbols {
		if v == -1 {
			st.symbols[k] = mem
			mem++
		}
	}
}
