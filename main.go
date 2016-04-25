package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ogier/pflag"
)

var (
	name    = "csvp"
	version = "0.10.0"

	flagset         = pflag.NewFlagSet(name, pflag.ContinueOnError)
	indexesList     = flagset.StringP("indexes", "i", "", "")
	headersList     = flagset.StringP("headers", "h", "", "")
	isTSV           = flagset.BoolP("tsv", "t", false, "")
	delimiter       = flagset.StringP("delimiter", "d", ",", "")
	outputDelimiter = flagset.StringP("output-delimiter", "D", "\t", "")
	isHelp          = flagset.BoolP("help", "", false, "")
	isVersion       = flagset.BoolP("version", "", false, "")
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

func printErr(err interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
}

func guideToHelp() {
	fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", name)
}

func toDelimiter(s string) (ch rune, err error) {
	s, err = strconv.Unquote(`"` + s + `"`)
	if err != nil {
		return 0, err
	}

	a := []rune(s)
	if len(a) != 1 {
		return 0, fmt.Errorf("the delimiter must be a single character")
	}
	return a[0], nil
}

func do(c *CSVScanner) error {
	for c.Scan() {
		fmt.Println(c.Text())
	}
	return c.Err()
}

func _main() int {
	flagset.SetOutput(ioutil.Discard)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	if *isHelp {
		printUsage()
		return 0
	}
	if *isVersion {
		printVersion()
		return 0
	}

	var selector Selector
	switch {
	case *indexesList != "" && *headersList != "":
		printErr("only one type of list may be specified")
		guideToHelp()
		return 2
	case *indexesList != "":
		selector = NewIndexes(*indexesList)
	case *headersList != "":
		selector = NewHeaders(*headersList)
	default:
		selector = NewAll()
	}

	c := NewCSVScanner(selector, nil)
	c.SetOutputDelimiter(*outputDelimiter)
	switch {
	case *isTSV:
		c.SetDelimiter('\t')
	default:
		ch, err := toDelimiter(*delimiter)
		if err != nil {
			printErr(err)
			guideToHelp()
			return 2
		}
		c.SetDelimiter(ch)
	}

	for _, file := range flagset.Args() {
		f, err := os.Open(file)
		if err != nil {
			printErr(err)
			return 1
		}
		defer f.Close()

		c.InitializeReader(f)
		if err := do(c); err != nil {
			printErr(err)
			return 1
		}
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
