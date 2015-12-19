csvp
====

[![Build Status](https://travis-ci.org/kusabashira/csvp.svg?branch=master)](https://travis-ci.org/kusabashira/csvp)

Print selected parts of CSV from each FILE to standard output.

```
$ cat items.csv
name,price,quantity
Apple,60,20
Grapes,140,8
Pineapple,400,2
Orange,50,14

$ cat items.csv | csvp -h price,quantity
60	20
140	8
400	2
50	14
```

Usage
-----

```
$ csvp [OPTION]... [FILE]...
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
```

Installation
------------

### go get

```
go get github.com/kusabashira/csvp
```

Options
-------

### --help

Display the usage and exit.

### --version

Output the version of csvp.

### -i, --indexes=LIST

Select only specified indexes.

Indexes separated by a `,`.

Each index starts from `1`, they are specified by the `index` or `range`.

```sh
# select only second column, and from fourth column to sixth column
csvp --indexes=2,4-6

# select only ninth column, seventh column, and up to third column
csvp --indexes=9,7,-3
```

#### index

`index` is a single index.

```sh
# select only second column
csvp --indexes=2

# select only forth column, first column, and second column
csvp --indexes=4,1,2
```

#### range

`range` are indexes from `first` to `last`.
It starts from the head if omitted `first`,
It continues until the end if omitted `last`.

```sh
# select only from the second column to the fourth column
csvp --indexes=2-4

# select only third column later
csvp --indexes=3-

# select only up to third column
csvp --indexes=-3

# select all columns
csvp --indexes=-
```

#### syntax of indexes list

Here is the syntax of indexes in extended BNF.

```
indexes = ( index | range ) , { "," , ( index | range ) }
range   = [ index ] , "-" , [ index ]
index   = { digit }
```

### -h, --headers=LIST

Select only specified headers.

Headers separated by a `,`.

```sh
# select only column of name
csvp --headers=name

# select only column of name, column of price, and column of quantity
csvp --indexes=name,price,quantity

# select only columns of "foo,bar" and columns of "baz"
csvp --indexes="foo\,bar,baz"
```

#### syntax of headers list

Here is the syntax of headers in extended BNF.

```
headers = header , { "," , header }
header  = { letter | "\," }
```

letter is a unicode character other than `,`.

### -t, -tsv

Change ths input delimiter to `\t`.  equivalent to -d'\t'.

### -d, --delimiter=DELIM

Change the input delimiter to `DELIM`.
`DELIM` is a unicode character.

```sh
# Read TSV
csvp --delimiter='\t'

# Read SSV
csvp --delimiter=' '
```

### -D, --output-delimiter=STRING

Change the output delimiter to `STRING`.
`STRING` is unicode characters.

```sh
# Outputs with a slash delimited
csvp --output-delimiter=/

# Outputs with a "::" delimited
csvp --output-delimiter=::
```

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
