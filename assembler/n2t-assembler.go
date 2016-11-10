package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/foggerty/flib"
	"github.com/foggerty/n2t/components"
)

var inputFile string
var outputFile string
var out *os.File

func main() {
	defineParams()

	AbortIf(
		func() bool { return strings.Trim(inputFile, "") != "" },
		func() { showHelp() })

	AbortIf(
		func() bool { return checkInput() },
		func() { fmt.Println("Cannot find input file.") })

	AbortIfErr(
		func() error { return setOutput() },
		"Error setting output.",
		nil)

	AbortIfErr(
		func() error { return Assemble(inputFile, out) },
		"Error when assembling.",
		func() { deleteOutput() })

	out.Sync()
	out.Close()
	os.Exit(0)
}

// Assemble will take a Hack assembler file (.asm) and writes to out,
// the "binary" machine codes.
func Assemble(in string, out *os.File) error {
	b, err := ioutil.ReadFile(in)

	if err != nil {
		return err
	}

	input := string(b)
	lexChan := components.StartLexingAsm(input)
	parser := components.NewParser(lexChan)

	if parser.Error != nil {
		return parser.Error
	}

	for {
		asm, ok := <-parser.Output

		if !ok {
			return parser.Error
		}

		fmt.Fprintln(out, asm)
	}
}

func defineParams() {
	flag.StringVar(&inputFile, "in", "", "Name of the input file.")
	flag.StringVar(&outputFile, "out", "",
		"Name of the output file (defaults to name of in, with the extension .hack).\n\nWill overwrite existing files.")

	flag.Parse()
}

func setOutput() error {

	if outputFile == "" {
		out = os.Stdout
		return nil
	}

	var err error
	out, err = os.Create(outputFile)

	return err
}

func checkInput() bool {
	_, err := os.Stat(inputFile)

	return err == nil
}

func showHelp() {
	fmt.Printf("\nNand2Tetris assembler.\n=====================\n\n")
	fmt.Printf("Usage:\n")

	flag.PrintDefaults()

	fmt.Println()
}

func deleteOutput() {
	if out != nil {
		if out.Close() != nil {
			panic("Something went horribly wrong trying to tidy")
		}
	}

	if _, exists := os.Stat(outputFile); !os.IsNotExist(exists) {
		// Note to self: find out how to flag warnings as "I know what I'm
		// doing".  The following line will only error with a *PathError,
		// which the above line has already taken care of.
		os.Remove(outputFile)
	}
}
