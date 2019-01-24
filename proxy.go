package gpac

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Proxy is proxy type defined in pac file
// like
// PROXY 127.0.0.1:8080
// SOCKS 127.0.0.1:1080
type Proxy struct {
	Type    string // Proxy type: PROXY HTTP HTTPS SOCKS DIRECT etc.
	Address string // Proxy address

	client *http.Client
	once   sync.Once
}

// IsDirect tests whether it is using direct connection
func (p *Proxy) IsDirect() bool {
	return p.Type == "DIRECT"
}

// URL returns a url representation for the proxy for curl -x
func (p *Proxy) URL() string {
	switch p.Type {
	case "DIRECT":
		return ""
	case "PROXY":
		return p.Address
	default:
		return fmt.Sprintf("%s://%s", strings.ToLower(p.Type), p.Address)
	}
}

// Proxy returns Proxy function that is ready use for http.Transport
func (p *Proxy) Proxy() func(*http.Request) (*url.URL, error) {
	var u *url.URL
	var err error

	switch p.Type {
	case "DIRECT":
		break
	case "PROXY":
		u, err = url.Parse(fmt.Sprintf("http://%s", p.Address))
	default:
		u, err = url.Parse(fmt.Sprintf("%s://%s", strings.ToLower(p.Type), p.Address))
	}

	return func(*http.Request) (*url.URL, error) {
		return u, err
	}
}

// Client returns an http.Client ready for use with this proxy
func (p *Proxy) Client() *http.Client {
	p.once.Do(func() {
		p.client = &http.Client{
			Transport: &http.Transport{
				Proxy: p.Proxy(),
			},
		}
	})
	return p.client
}

// Get issues a GET to the specified URL via the proxy
func (p *Proxy) Get(urlstr string) (*http.Response, error) {
	return p.Client().Get(urlstr)
}

// Do sends an HTTP request via the proxy and returns an HTTP response
func (p *Proxy) Do(req *http.Request) (*http.Response, error) {
	return p.Client().Do(req)
}

func (p *Proxy) String() string {
	return fmt.Sprintf("%s %s", p.Type, p.Address)
}

// ParseProxy parses proxy string returned by FindProxyForURL
// and returns a slice of proxies
func ParseProxy(pstr string) []*Proxy {
	var proxies []*Proxy
	ps := strings.FieldsFunc(pstr, func(r rune) bool {
		if r == ';' {
			return true
		}
		return false
	})

	for _, p := range ps {
		typeAddr := strings.Fields(p)
		if len(typeAddr) == 2 {
			proxies = append(proxies,
				&Proxy{
					Type:    strings.ToUpper(typeAddr[0]),
					Address: typeAddr[1],
				},
			)
		} else if len(typeAddr) == 1 {
			proxies = append(proxies,
				&Proxy{
					Type: strings.ToUpper(typeAddr[0]),
				},
			)
		}
	}

	return proxies
}
