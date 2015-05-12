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

func AverageTimespan(hardata *har.Har) time.Duration {
	var sum time.Duration
	entries := hardata.Log.Entries
	for i := 1; i < len(entries); i++ {
		sum += entries[i].Started.Sub(entries[i-1].Started)
	}
	return sum / time.Duration(len(entries)-1)
}
