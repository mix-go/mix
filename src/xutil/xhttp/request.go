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
			return doRequest(xReq, o)
		})
	}
	return doRequest(newXRequest(req), o)
}

func Do(req *http.Request, opts ...RequestOption) (*XResponse, error) {
	o := mergeOptions(opts)

	if o.RetryOptions != nil {
		xReq := newXRequest(req)
		return doRetry(o, func() (*XResponse, error) {
			xReq.RetryAttempts++
			return doRequest(xReq, o)
		})
	}
	return doRequest(newXRequest(req), o)
}

func doRequest(xReq *XRequest, opts *RequestOptions) (*XResponse, error) {
	var finalHandler HandlerFunc = func(xReq *XRequest, opts *RequestOptions) (*XResponse, error) {
		if !shutdownController.BeginRequest() {
			return nil, ErrShutdown
		}
		defer shutdownController.EndRequest()

		cli := http.Client{
			Timeout: opts.Timeout,
		}
		startTime := time.Now()
		r, err := cli.Do(xReq.Request)
		if err != nil {
			doDebug(opts, time.Now().Sub(startTime), xReq, nil, err)
			return nil, err
		}
		xResp := newXResponse(r)
		doDebug(opts, time.Now().Sub(startTime), xReq, xResp, nil)
		return xResp, nil
	}

	for i := len(opts.Middlewares) - 1; i >= 0; i-- {
		finalHandler = opts.Middlewares[i](finalHandler)
	}

	return finalHandler(xReq, opts)
}

func Shutdown() {
	shutdownController.InitiateShutdown()
}
