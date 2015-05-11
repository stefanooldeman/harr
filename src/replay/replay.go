package replay

import (
	"har"
	"net/http"
	"sync"
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

	epoch := entries[0].Started

	var done sync.WaitGroup
	done.Add(len(entries))

	client := new(http.Client)

	for _, entry := range entries {
		delay := entry.Started.Sub(epoch)
		go func(entry har.Entry) {
			time.Sleep(delay)
			err = Fire(client, &entry.Request)
			done.Done()
		}(entry)
	}

	done.Wait()
	return
}
