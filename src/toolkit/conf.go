package toolkit

type Config struct {
	Api     string
	Proxies []ProxyConfig
}

type ProxyConfig struct {
	Name     string
	Endpoint string
	To       string
}
