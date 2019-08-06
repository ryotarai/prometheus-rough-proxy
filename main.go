package main

import (
	"github.com/ryotarai/prometheus-rough-proxy/lib/cli"
	"log"
	"os"
)

func main() {
	if err := cli.Start(os.Args); err != nil {
		log.Fatal(err)
	}
}
