
# Book Search Demo [![Build Status](https://travis-ci.org/m-zajac/gobooksearchdemo.svg?branch=master)](https://travis-ci.org/m-zajac/gobooksearchdemo) [![Go Report Card](https://goreportcard.com/badge/github.com/m-zajac/gobooksearchdemo)](https://goreportcard.com/report/github.com/m-zajac/gobooksearchdemo)

A small project to demonstrate how I write go code today :)

## Goal

The goal is to write an appllication that fetches a book (from gutenber.org website), and then performs a fuzzy search given a phrase. Result should be the whole paragraph containg matching phrase. 
Application provides REST server and CLI version. Docker image is provided for easier use.

Code should be written as in serious project. Code readability, reliability are the most important. Code should be easy to test and expand.

## Description

### Search algorithm

Search is implemented using modified Levenstein distance algorithm.
For DP Levenstein distance algorithm see: https://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_full_matrix.
This code uses similar dynamic programming approach, but first dp array row is filled with zeroes.
Search returns first of best matches.

Time complexity is `O(nk)`, where `n` is the text size and `k` is search phraze size.

Space complexity is `O(n)`.

### Book source

Plain text version of books is fetched from gutenberg.org website. Once fetched, book is stored in local database (gzipped). Subsequent searches in a book will use local data.

## Finding id of a book

Go to gutenberg.org, find any book. You will see URL like: `http://gutenberg.org/ebooks/20`. "20" is the id of the book.

## Fun example

Search for a dragons in the bible :) `cli -id 10 -p 'dragon'`

```
Phrase found

32:32 For their vine is of the vine of Sodom, and of the fields of
Gomorrah: their grapes are grapes of gall, their clusters are bitter:
32:33 Their wine is the poison of dragons, and the cruel venom of
asps.
```

## Usage (docker)

Docker image is available on public docker hub.

Example usage:

### CLI

Example:

    docker run --rm -ti -v "$(pwd)"/data:/data mxzajac/gobooksearchdemo cli -l error -id 1513 -f 2 -p 'oh romeo romeo' 

Arguments:
- `-id` - book id
- `-p` - phrase to search (case insensitive)
- `-f` - fuzziness (maximum edit distance frract match to given phrase)
- `-l` - log level [debug, info, error (default)]
- `-dbf` - file for local books database (default: /data/bookdata.db)

### Server

Example:

    docker run --rm -ti -p 8080:8080 -v "$(pwd)"/data:/data mxzajac/gobooksearchdemo server

Arguments:
- `-p` - server port *(default: 8080)
- `-l` - log level [debug, info, error (default)]
- `-dbf` - file for local books database (default: /data/bookdata.db)

#### REST API Documetation

After running a server open `http://localhost:8080/docs` in a browser. You will see swagger UI. You can even perform some requests from this page.

Example curl:

    curl -X POST "http://localhost:8080/search" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"bookId\": \"1513\", \"fuzziness\": 0, \"phrase\": \"romeo\"}"