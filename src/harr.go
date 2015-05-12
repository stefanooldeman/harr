package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"har"
	"replay"
	"strconv"
)

func main() {
	usage := `Http ARchive Replayer

Usage:
  harr [options] HARFILE

Options:
  --target=URL, -t=URL               Target http url
  --external, -e                     Also run external requests
  --concurrency=NUMBER, -c=NUMBER    Number of concurrent requests to use
                                     [default: 1]
  --repeat=TIMES, -n=TIMES           Total number of replays to run
                                     [default: 1]
`

	arguments, _ := docopt.Parse(usage, nil, true, "harr 0.1-dev", false)
	arguments["--target"], _ = arguments["--target"].(string)

	fmt.Println(arguments)

	replayOptions := &replay.Options{
		Target: arguments["--target"].(string),
	}

	result := &har.Har{}
	if err := har.ParseFile(arguments["HARFILE"].(string), &result); err != nil {
		panic(err)
	}

	concurrency := 1
	if str, ok := arguments["--concurrency"].(string); ok {
		concurrency, _ = strconv.Atoi(str)
	}
	repeat := 1
	if str, ok := arguments["--repeat"].(string); ok {
		repeat, _ = strconv.Atoi(str)
	}
	if concurrency > repeat {
		concurrency = repeat
	}

	replayer := replay.NewAsyncReplayer(concurrency)
	defer replayer.Close()

	replayer.Replay(result, repeat, replayOptions)

	replayer.Wait()
	return
}
