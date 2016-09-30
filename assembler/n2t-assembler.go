package main

import (
	"flag"
	"fmt"
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

	err := components.Assemble(inputFile, out)

	if err != nil {
		dumpErr("Error when assembling.", err)
		deleteOutput()
		os.Exit(1)
	}

	out.Sync()
	out.Close()
	os.Exit(0)
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
