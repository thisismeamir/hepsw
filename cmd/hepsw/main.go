package main

import (
	_ "os"

	"github.com/thisismeamir/hepsw/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
