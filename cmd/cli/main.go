package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/m-zajac/gobooksearchdemo/internal/adapters/bookcache"
	"github.com/m-zajac/gobooksearchdemo/internal/adapters/gutenberg"
	"github.com/m-zajac/gobooksearchdemo/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	bookID := flag.String("id", "", "gutenberg.org book id")
	phrase := flag.String("p", "", "search phrase")
	fuziness := flag.Uint("f", 2, "fuzzinnes - maximum edit distance frract match to given phrase")
	logLevel := flag.String("l", "error", "log level [debug, info, error (default)]")
	dbFileName := flag.String("dbf", "./data/bookdata.db", "local db file name (default: /data/bookdata.db)")
	flag.Parse()

	if *bookID == "" {
		fmt.Println("book id is required")
		os.Exit(1)
	}
	if *phrase == "" {
		fmt.Println("phrase is required")
		os.Exit(1)
	}

	logrusLevel, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		fmt.Println("invalid log level")
		os.Exit(1)
	}
	logrus.SetLevel(logrusLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02T15:04:05.000",
	})

	bs := gutenberg.NewSource()

	bc, err := bookcache.NewBoltCache(*dbFileName)
	if err != nil {
		logrus.Fatalf("creating book cache: %v", err)
	}

	service, err := app.NewService(bs, bc)
	if err != nil {
		logrus.Fatalf("creating book search service: %v", err)
	}

	paragraph, err := service.FindBookParagraph(*bookID, *phrase, *fuziness)
	if err == app.ErrBookNotFound {
		fmt.Printf("Book with id '%s' not found\n", *bookID)
		return
	} else if err != nil {
		logrus.Fatalf("searching for paragraph in book '%s': %v", *bookID, err)
	}

	if paragraph == "" {
		fmt.Printf("Phrase not found.\n")
		return
	}
	fmt.Printf("Phrase found.\n\n%s\n", paragraph)
}
