package main

import (
	"helpers"
	"util"

	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
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
	// fmt.Println("======== REQ =========")
	// fmt.Println(string(b))
	// fmt.Println("======== RESP ========")
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
	//if ec.HasEvent(req) {
	//	// change resp if body exists
	//} else {
	//}

	//fmt.Println(string(b))
	return resp, nil
}

var DefaultTransport http.RoundTripper = &ProxyTransport{&http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
	// ExpectContinueTimeout: 2 * time.Second, // only go 1.6
}}

var ec *EventsContainer

func init() {
	ec = NewEventsContainer()
}

func main() {
	target, err := url.Parse("http://localhost:3002")
	if err != nil {
		panic(err)
	}
	targetQuery := target.RawQuery

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	proxy := &httputil.ReverseProxy{Director: director, Transport: DefaultTransport}
	http.Handle("/", proxy)

	go eventsApi()
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func eventsApi() {
	api := func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(ec.EventsList)
		if err != nil {
			fmt.Fprintf(w, "Error %v", err)
			return
		}
		fmt.Fprintf(w, "%v", string(b))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", api)

	assets := http.FileServer(&assetfs.AssetFS{Asset: util.Asset, AssetDir: util.AssetDir, AssetInfo: util.AssetInfo, Prefix: "src/util/assets"})
	mux.Handle("/", assets)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
