package toolkit

type Conf struct {
	Api     Endpoint
	Proxies Proxy
}

type Proxy struct {
	From Endpoint
	To   Endpoint
}

type Endpoint struct {
	Bind string
	Port string
}
