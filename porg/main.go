package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// bucket: porg.dev.brnstz.com
type config struct {
}

func main() {
	var (
		c   config
		err error
	)

	err = envconfig.Process("porg", &c)
	if err != nil {
		log.Fatal(err)
	}
}
