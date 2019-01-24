package gpac

import "strings"

// Proxy is proxy type defined in pac file
// like
// PROXY 127.0.0.1:8080
// SOCKS 127.0.0.1:1080
type Proxy struct {
	Type    string // Proxy type: PROXY HTTP HTTPS SOCKS DIRECT etc.
	Address string // Proxy address
}

// IsDirect test whether it is using direct connection
func (p *Proxy) IsDirect() bool {
	return p.Type == "DIRECT"
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
