package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/m-zajac/gobooksearchdemo/internal/adapters/bookcache"
	"github.com/m-zajac/gobooksearchdemo/internal/adapters/gutenberg"
	"github.com/m-zajac/gobooksearchdemo/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	bookID := flag.String("id", "", "gutenberg.org book id")
	phrase := flag.String("p", "", "search phrase")
	logLevel := flag.String("l", "error", "log level [debug, info, error (default)]")
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

	// bookID := "33598"                       //"10"
	// phrase := "słabe wymawianie samogłosek" // "ood and upright is the LORD: "

	bs := gutenberg.NewSource()

	bc, err := bookcache.NewBoltCache("./bookdata.db")
	if err != nil {
		log.Fatalf("creating book cache: %v", err)
	}

	service, err := app.NewService(bs, bc)
	if err != nil {
		log.Fatalf("creating book search service: %v", err)
	}

	paragraph, err := service.FindBookParagraph(*bookID, *phrase, 2)
	if err == app.ErrBookNotFound {
		fmt.Printf("Book with id '%s' not found\n", *bookID)
		return
	} else if err != nil {
		log.Fatalf("searching for paragraph in book '%s': %v", *bookID, err)
	}

	fmt.Printf("Found\n\n%s\n", paragraph)
}
