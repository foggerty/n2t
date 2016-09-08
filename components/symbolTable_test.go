package components

import (
	"testing"
)

// Stupid-simple tests to get used to the testing framework

func TestNewIsNotNill(t *testing.T) {
	test := NewSymbolTable()

	vNil := test.variables == nil
	lNil := test.labels == nil

	if vNil || lNil {
		t.Error("Variable or label map not initialised in new instance of SymbolTable.")
	}
}
