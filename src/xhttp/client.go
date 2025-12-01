package xhttp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrAbortRetry = errors.New("xhttp: abort further retries")

	ErrShutdown = errors.New("xhttp: service is currently being shutdown and will no longer accept new requests")
)

var DefaultClient = NewClient(&http.Client{}, DefaultOptions)

func Fetch(ctx context.Context, method string, u string, opts ...RequestOption) (*Response, error) {
	return DefaultClient.Fetch(ctx, method, u, opts...)
}

func NewRequest(method string, u string, opts ...RequestOption) (*Request, error) {
	return DefaultClient.NewRequest(method, u, opts...)
}

func Do(ctx context.Context, req *Request) (*Response, error) {
	return DefaultClient.Do(ctx, req)
}

func DoRequest(ctx context.Context, req *http.Request, opts ...RequestOption) (*Response, error) {
	return DefaultClient.DoRequest(ctx, req, opts...)
}

type Request struct {
	http.Request

	Body Body

	// Number of retries
	RetryAttempts int
}

type Response struct {
	http.Response

	Body Body
}

type Body []byte

func (t Body) String() string {
	return BytesToString(t)
}

type Client struct {
	http.Client
	DefaultOptions RequestOptions
}

func NewClient(c *http.Client, options RequestOptions) *Client {
	return &Client{
		Client:         *c,
		DefaultOptions: options,
	}
}

func newRequest(r *http.Request, reloading bool) *Request {
	if r == nil {
		return nil
	}
	req := &Request{
		Request: *r,
	}

	if r.Body == nil {
		return req
	}

	if reloading {
		// std body > xhttp body
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return req
		}
		req.Body = b
		r.Body = io.NopCloser(bytes.NewReader(b))
	}

	return req
}

func newResponse(r *http.Response) *Response {
	if r == nil {
		return nil
	}
	resp := &Response{
		Response: *r,
	}

	if r.Body == nil {
		return resp
	}

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return resp
	}
	resp.Body = b
	return resp
}

func (t *Client) Fetch(ctx context.Context, method string, u string, opts ...RequestOption) (*Response, error) {
	o := mergeOptions(t, opts)

	URL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: method,
		URL:    URL,
		Header: o.Header,
	}
	if o.Body != nil {
		// xhttp body > std body
		req.Body = io.NopCloser(bytes.NewReader(o.Body))
	}
	xReq := newRequest(req, false)

	if o.RetryOptions != nil {
		return doRetry(o, func() (*Response, error) {
			xReq.RetryAttempts++
			return t.doXRequest(ctx, xReq, o)
		})
	}
	return t.doXRequest(ctx, xReq, o)
}

func (t *Client) NewRequest(method string, u string, opts ...RequestOption) (*Request, error) {
	o := mergeOptions(t, opts)

	URL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: method,
		URL:    URL,
		Header: o.Header,
	}
	if o.Body != nil {
		// xhttp body > std body
		req.Body = io.NopCloser(bytes.NewReader(o.Body))
	}
	return newRequest(req, false), nil
}

func (t *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	o := mergeOptions(t, nil)

	if o.RetryOptions != nil {
		return doRetry(o, func() (*Response, error) {
			req.RetryAttempts++
			return t.doXRequest(ctx, req, o)
		})
	}

	return t.doXRequest(ctx, req, o)
}

func (t *Client) DoRequest(ctx context.Context, req *http.Request, opts ...RequestOption) (*Response, error) {
	o := mergeOptions(t, opts)
	xReq := newRequest(req, true)

	if o.Header != nil {
		xReq.Header = o.Header
	}

	if o.RetryOptions != nil {
		return doRetry(o, func() (*Response, error) {
			xReq.RetryAttempts++
			return t.doXRequest(ctx, xReq, o)
		})
	}

	return t.doXRequest(ctx, xReq, o)
}

func (t *Client) doXRequest(ctx context.Context, xReq *Request, opts *RequestOptions) (*Response, error) {
	var finalHandler HandlerFunc = func(xReq *Request, opts *RequestOptions) (*Response, error) {
		if !shutdownController.BeginRequest() {
			return nil, ErrShutdown
		}
		defer shutdownController.EndRequest()

		cli := t
		cli.Timeout = opts.Timeout
		startTime := time.Now()
		r, err := cli.Client.Do(&xReq.Request)
		if err != nil {
			t.doDebug(ctx, opts, time.Now().Sub(startTime), xReq, nil, err)
			return nil, err
		}
		xResp := newResponse(r)
		t.doDebug(ctx, opts, time.Now().Sub(startTime), xReq, xResp, nil)
		return xResp, nil
	}

	for i := len(opts.Middlewares) - 1; i >= 0; i-- {
		finalHandler = opts.Middlewares[i](finalHandler)
	}

	return finalHandler(xReq, opts)
}
