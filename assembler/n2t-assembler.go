package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/foggerty/n2t/components"
)

var inputFile string
var outputFile string
var out *os.File

func main() {
	defineParams()

	if strings.Trim(inputFile, "") == "" {
		showHelp()
		os.Exit(1)
	}

	if err := setOutput(); err != nil {
		dumpErr("Error setting output.", err)
		os.Exit(1)
	}

	if !checkInput() {
		os.Exit(1)
	}

	err := Assemble(inputFile, out)

	if err != nil {
		dumpErr("Error when assembling.", err)
		deleteOutput()
		os.Exit(1)
	}

	out.Sync()
	out.Close()
	os.Exit(0)
}

// Will take a Hack assembler file (.asm) and writes to out, the
// "binary" machine codes.
func Assemble(in string, out *os.File) error {
	b, err := ioutil.ReadFile(in)

	if err != nil {
		return err
	}

	input := string(b)

	// Create a Lexer, it will kick off the lexing in a go routine
	lexChan := components.StartLexingAsm(input)

	// Create a parser, and hand it the output from the lexer.  It will
	// run the first pass (building symbol table) and return any errors
	// found by the lexer.  i.e. the lexr and the first pass will be run
	// concurrently, but since we cannot move onto the second pass
	// (parsing the tokens and writing the file) until that's complete,
	// there's no benefit to running the second phase concurrently.
	parser, errs := components.NewParser(lexChan)

	if errs != nil {
		return errs.AsError()
	}

	// Note that the parser will stop writing after the first error that
	// it encounters.  It is up to the calling routine to tidy the file.
	return parser.Run(out).AsError()
}

func defineParams() {
	flag.StringVar(&inputFile, "in", "", "Name of the input file.")
	flag.StringVar(&outputFile, "out", "",
		"Name of the output file (defaults to name of in, with the extension .hack).  Will overwrite existing files.")

	flag.Parse()
}

func setOutput() error {

	if outputFile == "" {
		out = os.Stdout
		return nil
	}

	ext := filepath.Ext(outputFile)

	switch ext {
	case ".hack":
		break
	case "":
		outputFile += ".hack"
	default:
		var i int
		i = len(outputFile) - len(ext)
		outputFile = outputFile[:i] + ".hack"
	}

	var err error
	out, err = os.Create(outputFile)

	if err != nil {
		return err
	}

	return nil
}

func checkInput() bool {
	if _, err := os.Stat(inputFile); err == nil {
		return true
	}

	fmt.Println("Cannot find input file.")
	return false
}

func showHelp() {
	fmt.Printf("\nNand2Tetris assembler.\n=====================\n\n")
	fmt.Printf("Usage:\n")

	flag.PrintDefaults()

	fmt.Println()
}

func dumpErr(msg string, err error) {
	fmt.Printf("Something went horribly wrong:\n%s\n%s\n", msg, err.Error())
}

func deleteOutput() {
	if out != nil {
		out.Close()
	}

	if _, exists := os.Stat(outputFile); !os.IsNotExist(exists) {
		os.Remove(outputFile)
	}
}
