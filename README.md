A tool to watermark PDFs in a directory.

Uses [github.com/pdfcpu/pdfcpu](https://github.com/pdfcpu/pdfcpu)

## Version

1.0.1

## Usage

    mark file.pdf
    mark DIR // recursively scans DIR for PDFs 
    mark -d file|DIR // the -d flag causes mark to overwrite the original files.
    mark -s file|DIR // the -s flag prints a statistics file

## Install

A Windows installer is available on the [releases page](https://github.com/richardlehane/mark/releases).

For non-Windows platforms, build with [Go](https://golang.org).

## Rights

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)