package main

import (
	"flag"
	"fmt"
)

var apiPort = flag.String("api-port", "8080", "port for api")

// TODO: bind for api
var usage = flag.Bool("usage", false, "Show usage")

var proxysList proxys

func init() {
	proxysList = proxys(make([]string, 0))
	flag.Var(&proxysList, "proxy", "proxy definition like localhost:3000 localhost:3001")
}

type proxys []string

// this will be used as many times as flag was called
func (p proxys) Set(value string) error {
	p = append(p, value)
	fmt.Println(len(p))
	fmt.Println(value)
	return nil
}

func (p proxys) String() string {
	return fmt.Sprintf("%v", len(p))
}

func main() {
	flag.Parse()
	fmt.Println("aaa")
	fmt.Println(proxysList)
}
