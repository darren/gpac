package gpac_test

import (
	"testing"

	"github.com/darren/gpac"
)

func TestParseProxy(t *testing.T) {
	proxy := "PROXY 127.0.0.1:8080; SOCKs 127.0.0.1:1080; Direct"

	proxies := gpac.ParseProxy(proxy)

	if len(proxies) != 3 {
		t.Error("Parse failed")
		return
	}

	if proxies[1].Type != "SOCKS" {
		t.Error("Should be SOCKS5")
	}

	if !proxies[2].IsDirect() {
		t.Error("Should be direct")
	}
}
