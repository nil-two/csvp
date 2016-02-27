### v0.10.0 - 2016-02-27

- Ignore headers in all input-files if --headers specified.

### v0.9.0 - 2016-02-18

- Ignore backslash in front of character in headers.
  - `a\\b` is interpreted as `a\b`.
  - `a\bc\de` is interpreted as `abcde`.
  - Trailing backslash also ignored.

### v0.8.1 - 2016-01-12

- Release compiled binary for Windows, OSX, and Linux.

### v0.8.0 - 2015-12-13

- Allow mixed flag like "-ti 3,5".

### v0.7.0 - 2015-11-13

- Support -d, --delimiter to change input delimiter.
- Support -t, --tsv to read TSV.

### v0.6.0 - 2015-11-13

- Rename -d, --delimiter to -D, --output-delimiter.

### v0.5.0 - 2015-11-10

- Support range indexes like "3-5" and "4-".
  - "3-5" means "From 3 to 5".
  - "-3" means "From 1 to 3".
  - "3-" means "From 3 to last column".
  - "-" means "From 1 to last column".
- Allow empty headers.
- Allow empty indexes.
- Disallow 0 in indexes.
- Change the format of the version from "v0.5.0" to "0.5.0".

### v0.4.1 - 2015-10-23

- Prevent panic if non-existent index is specified.

### v0.4.0 - 2015-10-22

- Select all column if neither indexes nor headers are specified.

### v0.3.0 - 2015-10-22

- Support -d, --delimiter to change output delimiter.

### v0.2.0 - 2015-10-22

- Support -h, --headers to select column by headers.
- Start index from 1 in indexes.
- Allow non-existent index.

### v0.1.0 - 2015-10-21

- Initial release.
