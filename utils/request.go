package utils

import (
	"net/http"
	"io/ioutil"
	"bytes"
	"io"
)

type HttpResponse struct {
	Url string
	Response *http.Response
	Body []byte
	Err error
}

func Request(method string, url string, data io.Reader) *HttpResponse {
	ch := make(chan *HttpResponse)
	go func() {
		client := &http.Client{}

		req, err := http.NewRequest(method, url, data)
		if data != nil {
			req.Header.Add("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			ch <- &HttpResponse{ url, resp, nil, err }
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		ch <- &HttpResponse{ url, resp, body, err }
	}()
	return <-ch
}

func Get(url string) *HttpResponse {
	return Request(http.MethodGet, url, nil)
}

func Post(url string, data *bytes.Buffer) *HttpResponse {
	return Request(http.MethodPost, url, data)
}

