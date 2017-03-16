package main

import (
	"toolkit"
)

func main() {
	config := toolkit.Config{Api: "localhost:8080", Proxies: []toolkit.ProxyConfig{
		toolkit.ProxyConfig{"First", "localhost:4000", "http://localhost:5000"},
		toolkit.ProxyConfig{"Second", "localhost:4001", "http://localhost:5001"},
	},
	}
	toolkit.StartHTTP(config)
}
