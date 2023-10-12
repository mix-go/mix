package xhttp

import (
	"github.com/mix-go/xutil/xconv"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	resp := &XRequest{
		Request: r,
	}
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return resp
	}
	resp.Body = b
	return resp
}

func newXResponse(r *http.Response) *XResponse {
	if r == nil {
		return nil
	}
	resp := &XResponse{
		Response: r,
	}
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
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
		req.Body = io.NopCloser(strings.NewReader(o.Body.String()))
	}

	if o.RetryOptions != nil {
		return doRetry(o, func() (*XResponse, error) {
			return do(o, req)
		})
	}
	return do(o, req)
}

func Do(req *http.Request, opts ...RequestOption) (*XResponse, error) {
	o := mergeOptions(opts)

	if o.RetryOptions != nil {
		return doRetry(o, func() (*XResponse, error) {
			return do(o, req)
		})
	}
	return do(o, req)
}

func do(opts *requestOptions, req *http.Request) (*XResponse, error) {
	cli := http.Client{
		Timeout: opts.Timeout,
	}
	startTime := time.Now()
	r, err := cli.Do(req)
	if err != nil {
		doDebug(opts, time.Now().Sub(startTime), req, nil, err)
		return nil, err
	}
	resp := newXResponse(r)
	doDebug(opts, time.Now().Sub(startTime), req, resp, nil)
	return resp, nil
}
