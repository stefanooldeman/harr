package har

import (
	"time"
)

type Har struct {
	Log Log `json:"log"`
}

type Log struct {
	Version string  `json:"version"`
	Creator Creator `json:"creator"`
	Pages   []Page  `json:"pages"`
	Entries []Entry `json:"entries"`
}

type Creator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Page struct {
	Started     time.Time          `json:"startedDateTime"`
	Id          string             `json:"id"`
	Title       string             `json:"title"`
	PageTimings map[string]float32 `json:"pageTimings"`
}

type Entry struct {
	Started    time.Time          `json:"startedDateTime"`
	Time       float32            `json:"time"`
	Request    Request            `json:"request"`
	Response   Response           `json:"response"`
	Timings    map[string]float32 `json:"timings"`
	Connection string             `json:"connection"`
	PageRef    string             `json:"pageref"`
}

type Request struct {
	Method      string   `json:"method"`
	URL         string   `json:"url"`
	HttpVersion string   `json:"httpVersion"`
	Headers     []Header `json:"headers"`
	QueryString []string `json:"queryString"`
	Cookies     []Cookie `json:"cookies"`
	HeaderSize  uint64   `json:"headerSize"`
	BodySize    uint64   `json:"bodySize"`
}

type Response struct {
	Status       int      `json:"status"`
	StatusText   string   `json:"statusText"`
	HttpVersion  string   `json:"httpVersion"`
	Headers      []Header `json:"headers"`
	Cookies      []Cookie `json:"cookies"`
	Content      Content  `json:"content"`
	RedirectURL  string   `json:"redirectURL"`
	HeaderSize   uint64   `json:"headerSize"`
	BodySize     uint64   `json:"bodySize"`
	TransferSize uint64   `json:"_transferSize"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Expires  string `json:"expires"`
	HttpOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
}

type Content struct {
	Size        uint64 `json:"size"`
	MimeType    string `json:"mimeType"`
	Compression int    `json:"compression"`
}
