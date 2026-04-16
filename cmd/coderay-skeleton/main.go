package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bogdan-copocean/coderay-plugin/skeleton"
)

func main() {
	fileArg := flag.String("file", "", "path to source file; optional :START-END suffix")
	includeImports := flag.Bool("include-imports", false, "include import statements")
	symbol := flag.String("symbol", "", "filter to a class or top-level function (optional dotted path)")
	fileLineRange := flag.String("file-line-range", "", "START-END line range (mutually exclusive with :suffix on --file)")
	flag.Parse()

	if *fileArg == "" {
		fmt.Fprintln(os.Stderr, "usage: coderay-skeleton --file PATH [options]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	var flr *string
	if *fileLineRange != "" {
		flr = fileLineRange
	}

	out, err := skeleton.ReadFileSkeleton(*fileArg, skeleton.Options{
		IncludeImports: *includeImports,
		Symbol:         *symbol,
	}, flr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(out)
}
