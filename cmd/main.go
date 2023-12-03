package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/blck-snwmn/jsontomd"
)

func main() {
	var f string
	flag.StringVar(&f, "f", "", "json file path")
	flag.Parse()

	file, err := os.Open(f)
	if err != nil {
		panic(fmt.Sprintf("failed to open file: %+v", err))
	}
	decoder := json.NewDecoder(file)
	array, err := jsontomd.DecodeArray(decoder)
	if err != nil {
		panic(fmt.Sprintf("failed to decode: %+v", err))
	}
	md, err := jsontomd.EncodeMarkdown(array)
	if err != nil {
		panic(fmt.Sprintf("failed to encode: %+v", err))
	}
	fmt.Println(md)
}
