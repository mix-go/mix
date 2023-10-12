package xhttp

import (
	"bytes"
	"github.com/mix-go/xutil/xconv"
	"io"
	"net/http"
	"net/url"
	"time"
)

type XRequest struct {
	*http.Request
	Body Body
}

type XResponse struct {
	*http.Response
	Body Body
}

type Body []byte

func (t Body) String() string {
	return xconv.BytesToString(t)
}

func newXRequest(r *http.Request) *XRequest {
	if r == nil {
		return nil
	}
	req := &XRequest{
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

func newXResponse(r *http.Response) *XResponse {
	if r == nil {
		return nil
	}
	resp := &XResponse{
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

func Request(method string, u string, opts ...RequestOption) (*XResponse, error) {
	o := mergeOptions(opts)
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
		return doRetry(o, func() (*XResponse, error) {
			return doRequest(o, req)
		})
	}
	return doRequest(o, req)
}

func Do(req *http.Request, opts ...RequestOption) (*XResponse, error) {
	o := mergeOptions(opts)

	if o.RetryOptions != nil {
		return doRetry(o, func() (*XResponse, error) {
			return doRequest(o, req)
		})
	}
	return doRequest(o, req)
}

func doRequest(opts *requestOptions, req *http.Request) (*XResponse, error) {
	cli := http.Client{
		Timeout: opts.Timeout,
	}
	startTime := time.Now()
	xreq := newXRequest(req)
	r, err := cli.Do(req)
	if err != nil {
		doDebug(opts, time.Now().Sub(startTime), xreq, nil, err)
		return nil, err
	}
	xresp := newXResponse(r)
	doDebug(opts, time.Now().Sub(startTime), xreq, xresp, nil)
	return xresp, nil
}
