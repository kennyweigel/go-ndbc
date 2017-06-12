package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"

	"github.com/kennyweigel/go-ndbc/ndbc"
)

func main() {
	flagID := flag.String("id", "44030", "5 Digit Buoy ID")
	flagMax := flag.Int("max", 1, "Max number of records to return")

	flag.Parse()

	data := ndbc.GetBuoyData5Day(*flagID, *flagMax)

	var out bytes.Buffer
	json.Indent(&out, data, "=", "\t")
	out.WriteTo(os.Stdout)
}
