package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"toolkit"
)

var (
	cfgFile = flag.String("cfgFile", "cfg.json", "Json file with config")
)

func main() {
	flag.Parse()
	//config := toolkit.Config{Api: "localhost:8080", Proxies: []toolkit.ProxyConfig{
	//	toolkit.ProxyConfig{"First", "localhost:4000", "http://localhost:5000"},
	//	toolkit.ProxyConfig{"Second", "localhost:4001", "http://localhost:5001"},
	//},
	//}
	config := &toolkit.Config{}
	b, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	toolkit.StartHTTP(*config)
}
