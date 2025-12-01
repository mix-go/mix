package xhttp

import (
	"bytes"
	"errors"
	"github.com/mix-go/xutil/xconv"
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

func NewRequest(method string, u string, opts ...RequestOption) (*Response, error) {
	return DefaultClient.NewRequest(method, u, opts...)
}

type Request struct {
	*http.Request

	Body Body

	// Number of retries
	RetryAttempts int
}

type Response struct {
	*http.Response

	Body Body
}

type Body []byte

func (t Body) String() string {
	return xconv.BytesToString(t)
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

func newRequest(r *http.Request) *Request {
	if r == nil {
		return nil
	}
	req := &Request{
		Request: r,
	}

	if r.Body == nil {
		return req
	}

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return req
	}
	req.Body = b
	r.Body = io.NopCloser(bytes.NewReader(b))
	return req
}

func newResponse(r *http.Response) *Response {
	if r == nil {
		return nil
	}
	resp := &Response{
		Response: r,
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

func (t *Client) NewRequest(method string, u string, opts ...RequestOption) (*Response, error) {
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
		req.Body = io.NopCloser(bytes.NewReader(o.Body))
	}

	if o.RetryOptions != nil {
		xReq := newRequest(req)
		return doRetry(o, func() (*Response, error) {
			xReq.RetryAttempts++
			return t.doRequest(xReq, o)
		})
	}
	return t.doRequest(newRequest(req), o)
}

func (t *Client) SendRequest(req *http.Request, opts ...RequestOption) (*Response, error) {
	o := mergeOptions(t, opts)

	if o.RetryOptions != nil {
		xReq := newRequest(req)
		return doRetry(o, func() (*Response, error) {
			xReq.RetryAttempts++
			return t.doRequest(xReq, o)
		})
	}
	return t.doRequest(newRequest(req), o)
}

func (t *Client) doRequest(xReq *Request, opts *RequestOptions) (*Response, error) {
	var finalHandler HandlerFunc = func(xReq *Request, opts *RequestOptions) (*Response, error) {
		if !shutdownController.BeginRequest() {
			return nil, ErrShutdown
		}
		defer shutdownController.EndRequest()

		cli := t
		cli.Timeout = opts.Timeout
		startTime := time.Now()
		r, err := cli.Do(xReq.Request)
		if err != nil {
			t.doDebug(opts, time.Now().Sub(startTime), xReq, nil, err)
			return nil, err
		}
		xResp := newResponse(r)
		t.doDebug(opts, time.Now().Sub(startTime), xReq, xResp, nil)
		return xResp, nil
	}

	for i := len(opts.Middlewares) - 1; i >= 0; i-- {
		finalHandler = opts.Middlewares[i](finalHandler)
	}

	return finalHandler(xReq, opts)
}
