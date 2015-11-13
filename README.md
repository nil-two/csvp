csvp
====

Print selected parts of CSV from each FILE to standard output.

```
$ cat items.csv
name,price,quantity
Apple,60,20
Grapes,140,8
Pineapple,400,2
Orange,50,14

$ cat items.csv | csvp -h=price,quantity
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
  -D, --output-delimiter=STRING
                 use STRING as the output delimiter (default: \t)
  --help
                 display this help text and exit
  --version
                 output version information and exit
```

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
