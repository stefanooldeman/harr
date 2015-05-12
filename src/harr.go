package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"har"
	"replay"
)

func main() {
	usage := `Http ARchive Replayer

Usage:
  harr [options] HARFILE

Options:
  --target=URL, -t=URL          Target http url.
  --external, -e                Also run external requests.

`
	arguments, _ := docopt.Parse(usage, nil, true, "harr 0.1-dev", false)
	arguments["--target"], _ = arguments["--target"].(string)
	fmt.Println(arguments)

	result := &har.Har{}
	if err := har.ParseFile(arguments["HARFILE"].(string), &result); err != nil {
		panic(err)
	}

	replay.Replay(result, &replay.Options{
		Target: arguments["--target"].(string),
	})
}
