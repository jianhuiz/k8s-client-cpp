package main

import (
	"flag"
)

func init() {
	flag.Parse()
	flag.Lookup("v").Value.Set("6")
}

func main() { }