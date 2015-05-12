package replay

import (
	"har"
	"net/http"
	"time"
)

type Options struct {
	Target string
}

type HarError struct {
	error
}

func Fire(client *http.Client, request *har.Request) (err error) {
	var req *http.Request
	if req, err = request.ToHTTP(); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

func Replay(hardata *har.Har, opts *Options) (err error) {
	entries := hardata.Log.Entries

	if len(entries) == 0 {
		// Well, that was quick..
		return
	}

	client := new(http.Client)

	start := time.Now()
	epoch := entries[0].Started

	for _, entry := range entries {
		now := time.Now()
		delayFromStart := entry.Started.Sub(epoch)
		delay := delayFromStart - now.Sub(start)
		time.Sleep(delay)
		err = Fire(client, &entry.Request)
		// continue even when error
	}

	return
}
