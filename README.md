
# go booksearch demo [![Build Status](https://travis-ci.org/m-zajac/gobooksearchdemo.svg?branch=master)](https://travis-ci.org/m-zajac/gobooksearchdemo) [![Go Report Card](https://goreportcard.com/badge/github.com/m-zajac/gobooksearchdemo)](https://goreportcard.com/report/github.com/m-zajac/gobooksearchdemo) [![Docs](https://pkg.go.dev/badge/github.com/m-zajac/gobooksearchdemo)](https://pkg.go.dev/github.com/m-zajac/gobooksearchdemo)

TODO

### Docker

Docker image is available on public docker hub.

Example usage:

#### CLI:

Full example:

    docker run --rm -ti mxzajac/gobooksearchdemo cli -l error -id 1513 -f 2 -p 'oh romeo romeo' 

Arguments:
- '-id' - book id
- '-p' - phrase to search (case insensitive)
- '-f' - fuzziness (maximum edit distance frract match to given phrase)
- '-l' - log level [debug, info, error (default)]