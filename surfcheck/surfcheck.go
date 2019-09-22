package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/kennyweigel/go-ndbc/ndbc"
)

func main() {
	flagID := flag.String("id", "44030", "5 Digit Buoy ID")
	flagMax := flag.Int("max", 1, "Max number of records to return")

	flag.Parse()

	data := ndbc.GetBuoyData5Day(*flagID, *flagMax)

	jsonData, err := json.Marshal(map[string]interface{}{
		"data": data,
	})

	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, jsonData, "=", "\t")
	out.WriteTo(os.Stdout)
}
