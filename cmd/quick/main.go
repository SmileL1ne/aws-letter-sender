package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/Dias1c/aws-letter-sender/internal/desktop"
)

func main() {
	err := desktop.Run()
	switch {
	case errors.Is(err, desktop.ErrFlagsRequired):
		log.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	case err != nil:
		log.Fatal(err)
	default:
		log.Println("Successfully sent the messsages!")
	}
}
