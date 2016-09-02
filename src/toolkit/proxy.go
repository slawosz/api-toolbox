package toolkit

import (
	"helpers"

	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ProxyTransport struct {
	*http.Transport
}

func (p *ProxyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBytes, err := helpers.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}
	// here, do actual request
	resp, err := p.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	respBytes, err := helpers.DumpResponse(resp, false)
	if err != nil {
		panic(err)
	}
	respBodyBytes, err := helpers.DumpResponseBody(resp)
	if err != nil {
		panic(err)
	}
	e := &Event{URL: req.URL.String(), Method: req.Method, Req: reqBytes, Resp: respBytes, RespBody: respBodyBytes}
	e.SetUUID()
	ec.AddEvent(req, e)

	return resp, nil
}

func ProxyHandler(targetURL *url.URL) http.Handler {
	targetQuery := targetURL.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = singleJoiningSlash(targetURL.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director, Transport: getDefaultTransport()}
}

func getDefaultTransport() http.RoundTripper {
	return &ProxyTransport{&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   300 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		// ExpectContinueTimeout: 2 * time.Second, // only go 1.6
	}}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
