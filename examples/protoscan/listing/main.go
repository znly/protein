package main

import (
	"fmt"
	"log"

	"github.com/znly/protein/protoscan"
)

func main() {
	schemas, err := protoscan.ScanSchemas()
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range schemas {
		fmt.Printf("[%s] %s\n", s.GetUID(), s.GetFQName())
		for uid, name := range s.GetDeps() {
			fmt.Printf("\tdepends on: [%s] %s\n", uid, name)
		}
	}
}
