percent
=======

The tiny macro processor

```
$ type sample.in
Windows=%WINDOWS%

$ percent WINDOWS=C:\Windows sample.in
Windows=C:\Windows
```

Usage: `percent [-ansi] {NAME=VALUE}... filename...`

* `-ansi`
    * The encoding of VALUE is not UTF8 (ANSI)
