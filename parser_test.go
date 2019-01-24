package gpac_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/darren/gpac"
)

func Example() {
	pacf, _ := os.Open("testdata/wpad.dat")
	defer pacf.Close()

	data, _ := ioutil.ReadAll(pacf)
	pac, _ := gpac.New(string(data))

	r, _ := pac.FindProxyForURL("http://www.example.com/")

	fmt.Println(r)
	// Output:
	// PROXY 4.5.6.7:8080; PROXY 7.8.9.10:8080
}

func TestProxyGet(t *testing.T) {
	pacf, _ := os.Open("testdata/wpad.dat")
	defer pacf.Close()

	data, _ := ioutil.ReadAll(pacf)

	pac, err := gpac.New(string(data))
	if err != nil {
		t.Fatal(err)
	}

	proxies, err := pac.FindProxy("http://www.example.com/")
	if err != nil {
		t.Fatal(err)
	}

	if len(proxies) != 2 {
		t.Fatal("Find proxy failed")
	}

	if proxies[1].URL() != "7.8.9.10:8080" {
		t.Error("Get URL from proxy failed")
	}
}

// test start stared in proxy_test.go
func TestPacGet(t *testing.T) {
	pac, err := gpac.New(`
        function FindProxyForURL(url, host) {
            return "PROXY 127.0.0.1:9991; PROXY 127.0.0.1:9992; PROXY 127.0.0.1:8080; DIRECT"
        }
    `)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := pac.Get("http://localhost:8081")
	if err != nil {
		t.Fatal(err)
	}

	body := readBodyAndClose(resp)
	if body != "Example" {
		t.Errorf("Response not expected: %s", body)
	}
}
