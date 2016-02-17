package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ogier/pflag"
	"github.com/yuya-takeyama/argf"
)

var (
	name    = "csvp"
	version = "0.8.1"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `
Usage: %s [OPTION]... [FILE]...
Print selected parts of CSV from each FILE to standard output.

Options:
  -i, --indexes=LIST
                 select only these indexes
  -h, --headers=LIST
                 select only these headers
  -t, --tsv
                 equivalent to -d'\t'
  -d, --delimiter=DELIM
                 use DELIM instead of comma for field delimiter
  -D, --output-delimiter=STRING
                 use STRING as the output delimiter (default: \t)
  --help
                 display this help text and exit
  --version
                 output version information and exit
`[1:], name)
}

func printVersion() {
	fmt.Fprintln(os.Stderr, version)
}

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
}

type Option struct {
	IndexesList     string
	HeadersList     string
	IsTSV           bool
	Delimiter       string
	OutputDelimiter string
	IsHelp          bool
	IsVersion       bool
	Files           []string
}

func parseOption(args []string) (opt *Option, err error) {
	flag := pflag.NewFlagSet(name, pflag.ContinueOnError)
	flag.SetOutput(ioutil.Discard)

	opt = &Option{}
	flag.StringVarP(&opt.IndexesList, "indexes", "i", "", "")
	flag.StringVarP(&opt.HeadersList, "headers", "h", "", "")
	flag.BoolVarP(&opt.IsTSV, "tsv", "t", false, "")
	flag.StringVarP(&opt.Delimiter, "delimiter", "d", ",", "")
	flag.StringVarP(&opt.OutputDelimiter, "output-delimiter", "D", "\t", "")
	flag.BoolVarP(&opt.IsHelp, "help", "", false, "")
	flag.BoolVarP(&opt.IsVersion, "version", "", false, "")

	if err = flag.Parse(args); err != nil {
		return nil, err
	}

	opt.Files = flag.Args()
	return opt, nil
}

func toDelimiter(s string) (r rune, err error) {
	s, err = strconv.Unquote(`"` + s + `"`)
	if err != nil {
		return 0, err
	}

	runes := []rune(s)
	if len(runes) != 1 {
		return 0, fmt.Errorf("the delimiter must be a single character")
	}
	return runes[0], nil
}

func newCSVScannerFromOption(opt *Option) (c *CSVScanner, err error) {
	var selector Selector
	switch {
	case opt.IndexesList != "" && opt.HeadersList != "":
		return nil, fmt.Errorf("only one type of list may be specified")
	case opt.IndexesList != "":
		selector = NewIndexes(opt.IndexesList)
	case opt.HeadersList != "":
		selector = NewHeaders(opt.HeadersList)
	default:
		selector = NewAll()
	}

	reader, err := argf.From(opt.Files)
	if err != nil {
		return nil, err
	}

	c = NewCSVScanner(selector, reader)
	c.SetOutputDelimiter(opt.OutputDelimiter)
	switch {
	case opt.IsTSV:
		c.SetDelimiter('\t')
	default:
		delimiter, err := toDelimiter(opt.Delimiter)
		if err != nil {
			return nil, err
		}
		c.SetDelimiter(delimiter)
	}
	return c, nil
}

func do(c *CSVScanner) error {
	for c.Scan() {
		fmt.Println(c.Text())
	}
	return c.Err()
}

func guideToHelp() {
	fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", name)
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
		printUsage()
		return 0
	case opt.IsVersion:
		printVersion()
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
