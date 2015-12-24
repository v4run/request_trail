package main

import (
	"errors"
	"net/http"
)

type resp struct {
	code int
	url  string
}

type transportWrapper struct {
	*http.Transport
	redirectTrail []resp
}

// ErrFailedRequest is thrown when the status code of any response is greater than 399
var ErrFailedRequest = errors.New("Request failed")

// ErrOverRedirection is thrown when the number of redirection exceeds the allowed limit
var ErrOverRedirection = errors.New("Redirection limit exceeded")

func roundTrip(req *http.Request, t transportWrapper) ([]resp, *http.Response, error) {
	transport := t.Transport
	redirectTrail := t.redirectTrail
	if transport == nil {
		transport = http.DefaultTransport.(*http.Transport)
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return redirectTrail, nil, err
	}
	redirectTrail = append(redirectTrail, resp{code: res.StatusCode, url: req.URL.String()})
	if res != nil && res.StatusCode > 399 {
		return redirectTrail, nil, ErrFailedRequest
	}
	return redirectTrail, res, nil
}

func (t *transportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	redirectTrail, res, err := roundTrip(req, *t)
	t.redirectTrail = redirectTrail
	return res, err
}

func checkRedirect(maxRedirects int) func(*http.Request, []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if maxRedirects >= 0 && len(via) > maxRedirects {
			return ErrOverRedirection
		}
		return nil
	}
}
