package main

import (
	"github.com/joho/godotenv"
	"github.com/gopusher/gateway/comet"
	"flag"
	"github.com/gopusher/gateway/log"
)

func main() {
	filename := getArgs()

	log.Info("Load config file: %s", *filename)
	if err := godotenv.Load(*filename); err != nil {
		panic(err)
	}

	comet.Run()
}

func getArgs() (filename *string) {
	filename = flag.String("c", "./config.ini", "set config file path")
	flag.Parse()

	return
}
