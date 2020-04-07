package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type implHttpClient struct {
	hc         *http.Client
	beforeSend func(req *http.Request)
	afterRecv  func(resp *http.Response)
}

// Internal method to send request
// The returned body should be closed by invoker.
func (i *implHttpClient) send(url string, qs QueryString, form Form) (body io.ReadCloser, err error) {
	// append query string to URL
	if qs != nil {
		index := strings.IndexRune(url, '?')
		if index < 0 {
			url = fmt.Sprintf("%s?%s", url, qs.Encode())
		} else {
			url = fmt.Sprintf("%s&%s", url, qs.Encode())
		}
	}
	// process form
	method, data := http.MethodGet, io.Reader(nil)
	if form != nil {
		method, data = http.MethodPost, form.Archive()
	}
	// build request
	req, _ := http.NewRequest(method, url, data)
	if form != nil {
		req.Header.Set("Content-Type", form.ContentType())
	}
	if i.beforeSend != nil {
		i.beforeSend(req)
	}
	// send request
	if resp, err := i.hc.Do(req); err != nil {
		return nil, err
	} else {
		if i.afterRecv != nil {
			i.afterRecv(resp)
		}
		return resp.Body, nil
	}
}

func (i *implHttpClient) Get(url string, qs QueryString) ([]byte, error) {
	body, err := i.send(url, qs, nil)
	if err != nil {
		return nil, err
	}
	defer QuietlyClose(body)
	return ioutil.ReadAll(body)
}

func (i *implHttpClient) JsonApi(url string, qs QueryString, form Form, result interface{}) (err error) {
	body, err := i.send(url, qs, form)
	if err != nil {
		return
	}
	defer QuietlyClose(body)
	// parse response body
	d := json.NewDecoder(body)
	return d.Decode(result)
}

func (i *implHttpClient) JsonpApi(url string, qs QueryString, result interface{}) (err error) {
	body, err := i.send(url, qs, nil)
	if err != nil {
		return
	}
	defer QuietlyClose(body)
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	left, right := bytes.IndexByte(content, '('), bytes.LastIndexByte(content, ')')
	if left < 0 || right < 0 {
		return &json.SyntaxError{Offset: 0}
	}
	return json.Unmarshal(content[left+1:right], result)
}
