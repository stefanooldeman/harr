package replay

import (
	"har"
	"net/http"
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
