package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/yuya-takeyama/argf"
)

func usage() {
	os.Stderr.WriteString(`
Usage: csvp [OPTION]... [FILE]...
Print selected parts of CSV from each FILE to standard output.

Options:
  -i, --indexes=LIST       select only these indexes
  -h, --headers=LIST       select only these headers
  -d, --delimiter=STRING   use STRING as the output delimiter (default: \t)
      --help               display this help text and exit
      --version            display version information and exit
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.4.0
`[1:])
}

type Option struct {
	IndexesList string `short:"i" long:"indexes"`
	HeadersList string `short:"h" long:"headers"`
	Delimiter   string `short:"d" long:"delimiter" default:"\t"`
	IsHelp      bool   `          long:"help"`
	IsVersion   bool   `          long:"version"`
	Files       []string
}

func parseOption(args []string) (opt *Option, err error) {
	opt = &Option{}
	flag := flags.NewParser(opt, flags.PassDoubleDash)

	opt.Files, err = flag.ParseArgs(args)
	if err != nil && !opt.IsHelp && !opt.IsVersion {
		return nil, err
	}
	return opt, nil
}

func newCSVScannerFromOption(opt *Option) (c *CSVScanner, err error) {
	var selector Selector
	switch {
	case opt.IndexesList != "" && opt.HeadersList != "":
		return nil, fmt.Errorf("only one type of list may be specified")
	case opt.IndexesList != "":
		indexes, err := parseIndexesList(opt.IndexesList)
		if err != nil {
			return nil, err
		}
		selector = NewIndexes(indexes)
	case opt.HeadersList != "":
		headers, err := parseHeadersList(opt.HeadersList)
		if err != nil {
			return nil, err
		}
		selector = NewHeaders(headers)
	default:
		selector = NewAll()
	}

	reader, err := argf.From(opt.Files)
	if err != nil {
		return nil, err
	}

	c = NewCSVScanner(selector, reader)
	c.Delimiter = opt.Delimiter
	return c, nil
}

func do(c *CSVScanner) error {
	for c.Scan() {
		fmt.Println(c.Text())
	}
	return c.Err()
}

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "cspv:", err)
}

func guideToHelp() {
	os.Stderr.WriteString(`
Try 'cspv --help' for more information.
`[1:])
}

func _main() int {
	opt, err := parseOption(os.Args[1:])
	if err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	switch {
	case opt.IsHelp:
		usage()
		return 0
	case opt.IsVersion:
		version()
		return 0
	}

	c, err := newCSVScannerFromOption(opt)
	if err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	if err = do(c); err != nil {
		printErr(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
