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

// Slightly more sensible tests

func TestLabelIsStored(t *testing.T) {
	test := NewSymbolTable()
	test.AddLabel("Fred", 123)
	result, err := test.LabelLocation("Fred")

	if err != nil {
		t.Fatalf("Label lookup returned an error:\n%s", err.Error())
	}

	if result != 123 {
		t.Error("Label value was wrong")
	}
}

func TestDuplicateLabel(t *testing.T) {
	test := NewSymbolTable()
	test.AddLabel("Fred", 123)
	result := test.AddLabel("Fred", 456)

	if result == nil {
		t.Error("Was allowed to add same label twice.")
	}
}
