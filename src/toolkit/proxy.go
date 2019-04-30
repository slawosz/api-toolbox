package toolkit

import (
	"helpers"

	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func SetupProxy(proxies []ProxyConfig) {
	for _, p := range proxies {
		ec := NewEventsContainer()
		events[p.Name] = ec
		log.Printf("Proxying to %v", p.To)

		targetURL, err := url.Parse(p.To)
		if err != nil {
			panic(err)
		}
		r := http.NewServeMux()
		r.Handle("/", ProxyHandler(targetURL, ec))
		srv := &http.Server{
			Handler: r,
			Addr:    p.Endpoint,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Printf("Proxy is Listening on %v", p.Endpoint)
		go func() { log.Fatal(srv.ListenAndServe()) }()
	}
}

func ProxyHandler(targetURL *url.URL, ec *EventsContainer) http.Handler {
	targetQuery := targetURL.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		req.URL.Path = singleJoiningSlash(targetURL.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director, Transport: getDefaultTransport(ec)}
}

type ProxyTransport struct {
	*http.Transport
	ec *EventsContainer
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
	e := &Event{URL: req.URL.String(), Method: req.Method, Req: reqBytes, Resp: respBytes, RespBody: respBodyBytes, RespCode: resp.StatusCode}
	e.SetUUID()
	p.ec.AddEvent(req, e)

	return resp, nil
}

func getDefaultTransport(ec *EventsContainer) http.RoundTripper {
	return &ProxyTransport{
		&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   300 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			// ExpectContinueTimeout: 2 * time.Second, // only go 1.6
		},
		ec,
	}
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
