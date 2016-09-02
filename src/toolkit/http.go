package toolkit

import (
	"util"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

func StartHTTP() {
	go apiAndAssets()
	proxy()
}

func proxy() {
	targetURL, err := url.Parse("http://localhost:3002")
	if err != nil {
		panic(err)
	}

	http.Handle("/", ProxyHandler(targetURL))

	log.Fatal(http.ListenAndServe(":3001", nil))
}

func apiAndAssets() {
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
