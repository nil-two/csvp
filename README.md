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

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
