A script to watermark PDFs in a directory.

Uses [github.com/pdfcpu/pdfcpu](https://github.com/pdfcpu/pdfcpu)

## Version

1.0.0

## Usage

    mark file.pdf
    mark -d DIR // recursively scans DIR for PDFs. The -d flag causes it to overwrite the original files.

## Install

A Windows installer is available on the [releases page](https://github.com/richardlehane/mark/releases).

For non-Windows platforms, build with [Go](https://golang.org).

## Rights

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)