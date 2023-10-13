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

	// Number of retries
	RetryAttempts int
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
		xReq := newXRequest(req)
		return doRetry(o, func() (*XResponse, error) {
			xReq.RetryAttempts++
			return doRequest(o, xReq)
		})
	}
	return doRequest(o, newXRequest(req))
}

func Do(req *http.Request, opts ...RequestOption) (*XResponse, error) {
	o := mergeOptions(opts)

	if o.RetryOptions != nil {
		xReq := newXRequest(req)
		return doRetry(o, func() (*XResponse, error) {
			xReq.RetryAttempts++
			return doRequest(o, xReq)
		})
	}
	return doRequest(o, newXRequest(req))
}

func doRequest(opts *requestOptions, req *XRequest) (*XResponse, error) {
	cli := http.Client{
		Timeout: opts.Timeout,
	}
	startTime := time.Now()
	r, err := cli.Do(req.Request)
	if err != nil {
		doDebug(opts, time.Now().Sub(startTime), req, nil, err)
		return nil, err
	}
	xResp := newXResponse(r)
	doDebug(opts, time.Now().Sub(startTime), req, xResp, nil)
	return xResp, nil
}
