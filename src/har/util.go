package har

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func (self Request) ToHTTP() (req *http.Request, err error) {
	req, err = http.NewRequest(self.Method, self.URL, nil)
	req.Proto = self.HttpVersion

	for _, header := range self.Headers {
		req.Header[header.Name] = []string{header.Value}
	}
	for _, cookie := range self.Cookies {
		expires, _ := time.Parse(time.RFC3339, cookie.Expires)
		req.AddCookie(&http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Expires:  expires,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
		})
	}
	return
}

func Parse(dat []byte, res **Har) error {
	return json.Unmarshal(dat, res)
}

func ParseFile(file string, res **Har) error {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return Parse(raw, res)
}
