package requests

import (
	"io"
	"net/http"
)

type Response struct {
	Request       *Request
	RawResponse   *http.Response
	StatusCode    int
	Text          string
	Content       []byte
	Headers       Headers
	Cookies       []*http.Cookie
	ContentLength int64
}

func (h *Headers) Get(key string) string {
	v := (*h)[key]
	switch v.(type) {
	case string:
		return v.(string)
	case []string:
		return (v.([]string))[0]
	default:
		return ""
	}
}

func NewResponse(res *http.Response) (*Response, error) {
	response := &Response{}
	response.Headers = make(Headers)
	response.RawResponse = res
	response.ContentLength = res.ContentLength
	for k, v := range res.Header {
		response.Headers[k] = v
	}
	response.Cookies = res.Cookies()
	response.StatusCode = res.StatusCode
	var err error
	response.Content, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response.Text = string(response.Content)
	return response, nil
}
