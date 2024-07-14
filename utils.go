package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"
)

func NewRequest(opts []interface{}) (*Request, error) {
	request := &Request{}
	// 默认开启自动重定向
	request.AllowRedirects = true
	for _, opt := range opts {
		switch opt.(type) {
		case string:
			if strings.HasPrefix(opt.(string), "http://") || strings.HasPrefix(opt.(string), "https://") {
				request.Url = opt.(string)
			} else {
				if opt.(string) == "" {
					return nil, errors.New("method can not be empty")
				}
				request.Method = opt.(string)
			}
		case Headers:
			request.Headers = opt.(Headers)
		case Data:
			request.Data = opt.(Data)
		case Json:
			request.Json = opt.(Json)
		case Params:
			request.params = opt.(Params)
		case Proxy:
			if opt.(Proxy)["http"] == "" && opt.(Proxy)["https"] == "" {
				return nil, errors.New("proxy can not be empty")
			}
			request.Proxy = opt.(Proxy)
		case []FileRequest:
			request.Files = opt.([]FileRequest)
		case *FileLocal:
			request.FileLocal = opt.(*FileLocal)
		case *FileBytes:
			request.FileBytes = opt.(*FileBytes)
		case int:
			request.Timeout = opt.(int)
		case bool:
			request.AllowRedirects = opt.(bool)
		case RawData:
			request.RawData = opt.(RawData)
		case []byte:
			request.RawData = opt.([]byte)
		default:
			panic("Unsupported option")
		}
	}
	return request, nil
}

func DoRequest(r *Request, c *http.Client) (*Response, error) {
	var req *http.Request
	params, err := handleParams(r.params)
	if err != nil {
		return nil, err
	}
	query := params.Encode()
	if query != "" {
		r.Url = fmt.Sprintf("%s?%s", r.Url, query)
	}
	h, err := handleHeaders(r.Headers)
	if err != nil {
		return nil, err
	}
	form, err := handleData(r.Data)
	if err != nil {
		return nil, err
	}
	if r.Timeout > 0 {
		c.Timeout = time.Duration(r.Timeout) * time.Millisecond
	}
	reqBody := &bytes.Buffer{}
	var writer *multipart.Writer
	//多文件上传处理
	if r.Files != nil {
		writer = multipart.NewWriter(reqBody)
		for _, file := range r.Files {
			err := handleFile(writer, file)
			if err != nil {
				return nil, err
			}
		}
	}
	if r.FileBytes != nil {
		if writer == nil {
			writer = multipart.NewWriter(reqBody)
		}
		err := handleFile(writer, r.FileBytes)
		if err != nil {
			return nil, err
		}
	}
	if r.FileLocal != nil {
		if writer == nil {
			writer = multipart.NewWriter(reqBody)
		}
		err := handleFile(writer, r.FileLocal)
		if err != nil {
			return nil, err
		}
	}
	if writer != nil {
		// 添加其他表单字段
		for key, values := range form {
			for _, value := range values {
				err := writer.WriteField(key, value)
				if err != nil {
					return nil, err
				}
			}
		}
		writer.Close()
		req, err = http.NewRequest(r.Method, r.Url, reqBody)
		if h.Get("Content-Type") == "" {
			h.Set("Content-Type", writer.FormDataContentType())
		}
	} else if len(form) != 0 {
		req, err = http.NewRequest(r.Method, r.Url, strings.NewReader(form.Encode()))
		if h.Get("Content-Type") == "" {
			h.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	} else if r.Json != nil {
		jsonStr, err := json.Marshal(r.Json)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(r.Method, r.Url, strings.NewReader(string(jsonStr)))
		if h.Get("Content-Type") == "" {
			h.Set("Content-Type", "application/json")
		}
	} else if r.RawData != nil {
		req, err = http.NewRequest(r.Method, r.Url, bytes.NewReader(r.RawData))
	} else {
		req, err = http.NewRequest(r.Method, r.Url, nil)
	}
	if h != nil {
		req.Header = h
	}
	// 处理代理
	if r.Proxy != nil {
		c.Transport = &http.Transport{Proxy: func(rawR *http.Request) (*url.URL, error) {
			var proxyUrl *url.URL
			if v, ok := r.Proxy["socks5"]; ok {
				proxyUrl, err = url.Parse(v)
				if err != nil {
					return nil, err
				}
			} else {
				proxyUrl, err = url.Parse(r.Proxy[rawR.URL.Scheme])
			}
			if err != nil {
				return nil, err
			}
			return proxyUrl, nil
		}}
	}
	if !r.AllowRedirects {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	httpRes, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	response, err := NewResponse(httpRes)
	if err != nil {
		return nil, err
	}
	response.Request = r
	return response, nil
}

func handleHeaders(headers Headers) (http.Header, error) {
	h := http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"},
		"Connection": {"close"},
	}
	if headers == nil {
		return h, nil
	}
	for k, v := range headers {
		switch v.(type) {
		case string:
			h.Set(k, v.(string))
		case []string:
			h[k] = v.([]string)
		default:
			return nil, errors.New("invalid header type")
		}
	}
	return h, nil
}

func handleParams(params Params) (url.Values, error) {
	if params == nil {
		return nil, nil
	}
	p := url.Values{}
	for k, v := range params {
		switch v.(type) {
		case string:
			p.Set(k, v.(string))
		case []string:
			p[k] = v.([]string)
		default:
			return nil, errors.New("invalid param type")
		}
	}
	return p, nil
}

func handleData(d Data) (url.Values, error) {
	if d == nil {
		return nil, nil
	}
	form := url.Values{}
	if d == nil {
		return nil, nil
	}
	for k, v := range d {
		switch v.(type) {
		case string:
			form.Set(k, v.(string))
		case []string:
			form[k] = v.([]string)
		default:
			return nil, errors.New("invalid data type")
		}
	}
	return form, nil
}

// 处理文件上传情况
func handleFile(w *multipart.Writer, f FileRequest) error {

	h := make(textproto.MIMEHeader)
	r, err := f.GetReader()
	if err != nil {
		return err
	}
	defer r.Close()
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			f.GetFieldName(), f.GetFileName()))

	if f.GetContentType() != "" {
		h.Set("Content-Type", f.GetContentType())
	} else {
		// 文件MIME判断
		content, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			return err
		}
		h.Set("Content-Type", http.DetectContentType(content))
		r = io.NopCloser(bytes.NewReader(content))
	}
	part, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	io.Copy(part, r)
	return nil
}

func UrlEncodeFully(str string) string {
	s := ""
	for _, v := range str {
		s += fmt.Sprintf("%%%x", v)
	}
	return s
}

func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

func UrlDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}
