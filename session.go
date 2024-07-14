package requests

import (
	"net/http"
	"net/http/cookiejar"
)

type Session struct {
	Jar    http.CookieJar
	Client *http.Client
}

func NewSession() (*Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c := &http.Client{Jar: jar}
	return &Session{Jar: jar, Client: c}, nil
}

func (s *Session) Get(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "GET"
	if err != nil {
		return nil, err
	}
	response, err := DoRequest(request, s.Client)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (s *Session) Post(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "POST"
	if err != nil {
		return nil, err
	}
	response, err := DoRequest(request, s.Client)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Session) Put(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "PUT"
	if err != nil {
		return nil, err
	}
	response, err := DoRequest(request, s.Client)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Session) Delete(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "DELETE"
	if err != nil {
		return nil, err
	}
	response, err := DoRequest(request, s.Client)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Session) Options(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "OPTIONS"
	if err != nil {
		return nil, err
	}
	response, err := DoRequest(request, s.Client)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Session) Head(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	request.Method = "HEAD"
	if err != nil {
		return nil, err
	}
	response, err := DoRequest(request, s.Client)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Session) RequestCustom(opts ...any) (*Response, error) {
	request, err := NewRequest(opts)
	if err != nil {
		return nil, err
	}
	response, err := DoRequest(request, s.Client)
	if err != nil {
		return nil, err
	}
	return response, nil
}
