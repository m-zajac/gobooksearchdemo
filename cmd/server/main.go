package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/m-zajac/gobooksearchdemo/internal/adapters/bookcache"
	"github.com/m-zajac/gobooksearchdemo/internal/adapters/gutenberg"
	"github.com/m-zajac/gobooksearchdemo/internal/api"
	"github.com/m-zajac/gobooksearchdemo/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	port := flag.Uint("p", 8080, "server port")
	logLevel := flag.String("l", "info", "log level [debug, info (default), error]")
	dbFileName := flag.String("dbf", "./data/bookdata.db", "local db file name (default: ./data/bookdata.db)")
	flag.Parse()

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

	mux := api.NewMux(service)
	server := api.NewServer(
		*port,
		mux,
	)
	server.Run()
}
