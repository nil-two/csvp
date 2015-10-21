package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/yuya-takeyama/argf"
)

func usage() {
	os.Stderr.WriteString(`
Usage: cspv OPTION... [FILE]...
Print selected parts of CSV from each FILE to standard output.

Options:
  -i, --indexes=LIST       select only these indexes
      --help               display this help text and exit
      --version            display version information and exit
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.1.0
`[1:])
}

type Option struct {
	List      string `short:"i" long:"indexes"`
	IsHelp    bool   `          long:"help"`
	IsVersion bool   `          long:"version"`
	Files     []string
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
	indexes, err := parseIndexesList(opt.List)
	if err != nil {
		return nil, err
	}
	reader, err := argf.From(opt.Files)
	if err != nil {
		return nil, err
	}
	return NewCSVScanner(indexes, reader), nil
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
