package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Headers map[string]interface{}

type Params map[string]interface{}

type Data map[string]interface{}

type RawData []byte

type Json map[string]interface{}

type Proxy map[string]string

type FileLocal struct {
	Name        string
	Filename    string
	LocalPath   string
	ContentType string
}

type FileBytes struct {
	Name        string
	Filename    string
	Data        []byte
	ContentType string
}

type FileRequest interface {
	GetReader() (io.ReadCloser, error)
	GetFieldName() string
	GetFileName() string
	GetContentType() string
}

type Request struct {
	RawRequest *http.Request
	Method     string
	Headers    Headers
	Url        string
	params     Params
	Data       Data
	RawData    RawData
	Json       Json
	// 是否自动重定向， 默认true
	AllowRedirects bool
	Proxy          Proxy
	Timeout        int
	Files          []FileRequest
	FileLocal      *FileLocal
	FileBytes      *FileBytes
}

type ReqI interface {
	Get() (*Response, error)
	Post() (*Response, error)
	Put() (*Response, error)
	Delete() (*Response, error)
	Option() (*Response, error)
	Head() (*Response, error)
	RequestCustom() (*Response, error)
}

func Get(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "GET"
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	response, err := DoRequest(request, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func Post(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "POST"
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	response, err := DoRequest(request, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func Put(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "PUT"
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	response, err := DoRequest(request, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func Delete(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "DELETE"
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	response, err := DoRequest(request, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func Options(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "OPTIONS"
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	response, err := DoRequest(request, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func Head(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "HEAD"
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	response, err := DoRequest(request, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func RequestCustom(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	response, err := DoRequest(request, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *Request) Get() (*Response, error) {
	c := &http.Client{}
	r.Method = "GET"
	response, err := DoRequest(r, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *Request) Post() (*Response, error) {
	c := &http.Client{}
	r.Method = "POST"
	response, err := DoRequest(r, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *Request) Put() (*Response, error) {
	c := &http.Client{}
	r.Method = "PUT"
	response, err := DoRequest(r, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (r *Request) Delete() (*Response, error) {
	c := &http.Client{}
	r.Method = "DELETE"
	response, err := DoRequest(r, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (r *Request) Option() (*Response, error) {
	c := &http.Client{}
	r.Method = "OPTION"
	response, err := DoRequest(r, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (r *Request) Head() (*Response, error) {
	c := &http.Client{}
	r.Method = "HEAD"
	response, err := DoRequest(r, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *Request) RequestCustom() (*Response, error) {
	c := &http.Client{}
	response, err := DoRequest(r, c)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (f *FileLocal) GetReader() (io.ReadCloser, error) {
	file, err := os.Open(f.LocalPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileBytes) GetReader() (io.ReadCloser, error) {
	buf := bytes.NewReader(f.Data)
	return io.NopCloser(buf), nil
}

func (f *FileLocal) GetFieldName() string {
	if f.Name == "" {
		f.Name = "file"
	}
	return f.Name
}

func (f *FileBytes) GetFieldName() string {
	if f.Name == "" {
		f.Name = "file"
	}
	return f.Name
}

func (f *FileBytes) GetFileName() string {
	return f.Filename
}

func (f *FileLocal) GetFileName() string {
	if f.Filename == "" {
		f.Filename = filepath.Base(f.LocalPath)
	}
	return f.Filename
}

func (f *FileBytes) GetContentType() string {
	return f.ContentType
}

func (f *FileLocal) GetContentType() string {
	return f.ContentType
}

func (h Headers) String() string {

	formatHeader := make(map[string]string, len(h))
	for k, v := range h {
		formatHeader[k] = v.([]string)[0]
	}
	jsonStr, err := json.MarshalIndent(formatHeader, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonStr)
}
