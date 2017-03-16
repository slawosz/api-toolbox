package toolkit

import (
	"util"

	"encoding/json"
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	"github.com/kr/pretty"
)

var events map[string]*EventsContainer

func init() {
	events = make(map[string]*EventsContainer)
}

func StartHTTP(conf Config) {
	SetupProxy(conf.Proxies)
	apiAndAssets(conf.Api)
}

func apiAndAssets(apiURL string) {
	api := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		el := events[vars["proxy"]]
		b, err := json.Marshal(el)
		if err != nil {
			pretty.Println(el)
			fmt.Fprintf(w, "Error in EventContainer marshaling: %v", err)
			return
		}
		fmt.Fprintf(w, "%v", string(b))
	}

	proxies := func(w http.ResponseWriter, r *http.Request) {
		var proxies []string
		for p, _ := range events {
			proxies = append(proxies, p)
		}
		b, err := json.Marshal(proxies)
		if err != nil {
			fmt.Fprintf(w, "Error %v", err)
			return
		}
		fmt.Fprintf(w, "%v", string(b))
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/{proxy}", api)
	r.HandleFunc("/proxies", proxies)

	assets := http.FileServer(&assetfs.AssetFS{Asset: util.Asset, AssetDir: util.AssetDir, AssetInfo: util.AssetInfo, Prefix: "src/util/assets"})
	r.PathPrefix("/").Handler(assets)

	log.Printf("Webapp on %v", apiURL)
	srv := &http.Server{
		Handler: r,
		Addr:    apiURL,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
