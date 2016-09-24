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

// HACK - internally stored as ints so I can use -1 as a flag value.
type symbolTable struct {
	symbols     map[string]int
	initialised bool
	mem         asm
}

func newSymbolTable() symbolTable {
	return symbolTable{
		symbols: make(map[string]int),
		mem:     16,
	}
}

// Because we don't know the order in which variables and labels will
// be added - parser could see @START......(START) - we're just
// writing a flag value for all variables, in case they turn out to be
// labels.  Easier than updating them as we go and then reshuffling
// the variable locations.
func (st *symbolTable) witeMem() {
	mem := 16

	for k, v := range st.symbols {
		if v == -1 {
			st.symbols[k] = mem
			mem++
		}
	}
}

// Will always know the value of a label
func (st *symbolTable) addLabel(s string, m asm) {
	st.symbols[s] = int(m)
}

// We may be adding a variable, in which case set it to -1
// But if it already exists in the table, must be a label, so ignore.
func (st *symbolTable) addVariable(s string) {
	if _, ok := st.symbols[s]; !ok {
		st.symbols[s] = -1
	}
}

func (st *symbolTable) symbolValue(s string) (asm, error) {
	if res, ok := st.symbols[s]; ok {
		return asm(res), nil
	}

	msg := fmt.Sprintf("Internal error - %s not found in the lookup table.", s)

	return asm(0), errors.New(msg)
}
