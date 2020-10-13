package rhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiClient struct {
	Url     string
	Timeout time.Duration
	Header  http.Header
}

func (c *ApiClient) getUrl(uri string) string {
	return c.Url + uri
}

func (c *ApiClient) byteHttpClient(method string, uri string, jsonBody map[string]interface{}, urlValues url.Values, header http.Header) (ResponseByte, error) {
	var r ResponseByte

	var body io.Reader

	if jsonBody != nil {
		byteBody, err := json.Marshal(jsonBody)
		if err != nil {
			return r, err
		}
		body = bytes.NewBuffer(byteBody)
	} else {
		body = strings.NewReader(urlValues.Encode())
	}

	r.Url = c.getUrl(uri)
	req, err := http.NewRequest(method, r.Url, body)
	if err != nil {
		return r, err
	}
	//req.Header = http.Header{}

	if c.Header != nil {
		req.Header = c.Header.Clone()
	}

	for i, values := range header {
		for _, v := range values {
			req.Header.Add(i, v)
		}
	}

	if c.Timeout == 0 {
		c.Timeout = 10 * time.Second
	}

	client := &http.Client{
		Timeout: c.Timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}
	r.Code = resp.StatusCode
	defer resp.Body.Close()
	r.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	return r, nil
}

func (c *ApiClient) RowHttpClient(method string, uri string, jsonBody map[string]interface{}, urlValues url.Values, header http.Header) (ResponseRaw, error) {
	var r ResponseRaw

	bodyBytes, err := c.byteHttpClient(method, uri, jsonBody, urlValues, header)
	if err != nil {
		return r, err
	}

	r = ResponseRaw{
		Response: r.Response,
		Body:     string(bodyBytes.Body),
	}

	return r, nil
}

func (c *ApiClient) HttpClient(method string, uri string, jsonBody map[string]interface{}, urlValues url.Values, object interface{}, header http.Header) (ResponseJson, error) {
	var r ResponseJson

	if header == nil {
		header = http.Header{}
	}

	bodyBytes, err := c.byteHttpClient(method, uri, jsonBody, urlValues, header)
	if err != nil {
		return r, err
	}

	r.Raw = string(bodyBytes.Body)
	r.Response = bodyBytes.Response
	err = json.Unmarshal(bodyBytes.Body, &object)
	r.Body = &object
	if err != nil {
		return r, err
	}
	return r, nil
}
