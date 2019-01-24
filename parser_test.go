package gpac_test

import (
	"fmt"
	"io/ioutil"
	"os"

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
