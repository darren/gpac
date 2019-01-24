package gpac

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/robertkrimen/otto"
)

// Parser the parsed pac instance
type Parser struct {
	vm *otto.Otto
}

// FindProxyForURL  finding proxy for url
func (p *Parser) FindProxyForURL(urlstr string) (string, error) {
	u, err := url.Parse(urlstr)
	if err != nil {
		return "", err
	}

	f := fmt.Sprintf("FindProxyForURL('%s', '%s')", urlstr, u.Hostname())
	r, err := p.vm.Run(f)
	if err != nil {
		return "", err
	}
	return r.String(), nil
}

// New create a parser from text content
func New(text string) (*Parser, error) {
	vm := otto.New()
	registerBuiltinNatives(vm)
	registerBuiltinJS(vm)

	_, err := vm.Run(text)
	if err != nil {
		return nil, err
	}

	return &Parser{vm}, nil
}

func registerBuiltinJS(vm *otto.Otto) {
	_, err := vm.Run(builtinJS)
	if err != nil {
		panic(err)
	}
}

func registerBuiltinNatives(vm *otto.Otto) {
	for name, function := range builtinNatives {
		err := vm.Set(name, function)
		if err != nil {
			panic(err)
		}
	}
}

func fromReader(r io.ReadCloser) (*Parser, error) {
	defer r.Close()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return New(string(buf))
}

// FromFile load pac from file
func FromFile(filename string) (*Parser, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return fromReader(f)
}

// FromURL load pac from url
func FromURL(urlstr string) (*Parser, error) {
	resp, err := http.Get(urlstr)
	if err != nil {
		return nil, err
	}
	return fromReader(resp.Body)
}
