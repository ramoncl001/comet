package rest

import (
	"context"
	"net/url"
)

type RequestMethod string

func (r RequestMethod) String() string {
	return string(r)
}

func (r RequestMethod) Method() string {
	switch r {
	case LIST:
		return "GET"
	default:
		return string(r)
	}
}

const (
	LIST   RequestMethod = "LIST"
	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	PUT    RequestMethod = "PUT"
	DELETE RequestMethod = "DELETE"
	PATCH  RequestMethod = "PATCH"
)

type Request struct {
	Url           *url.URL
	Method        string
	QueryParams   map[string][]string
	PathParams    map[string]string
	Headers       map[string][]string
	Body          []byte
	UserAgent     string
	RemoteAddress string
	ctx           context.Context
}

func (r *Request) Context() context.Context {
	return r.ctx
}

func (r *Request) WithContext(ctx context.Context) *Request {
	return &Request{
		ctx:         ctx,
		Url:         r.Url,
		Method:      r.Method,
		QueryParams: r.QueryParams,
		PathParams:  r.PathParams,
		Headers:     r.Headers,
		Body:        r.Body,
	}
}
