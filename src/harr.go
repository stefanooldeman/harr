package main

import (
	"har"
	"os"
	"replay"
)

func main() {
	result := &har.Har{}
	if err := har.ParseFile(os.Args[1], &result); err != nil {
		panic(err)
	}

	replay.Replay(result)
}
