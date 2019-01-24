package gpac

import (
	"net"

	"github.com/robertkrimen/otto"
)

var builtinNatives = map[string]func(call otto.FunctionCall) otto.Value{
	"dnsResolve":  dnsResolve,
	"myIpAddress": myIPAddress,
}

func dnsResolve(call otto.FunctionCall) otto.Value {
	arg := call.Argument(0)
	if arg.IsNull() || arg.IsUndefined() {
		return otto.NullValue()
	}

	host := arg.String()
	ips, err := net.LookupIP(host)
	if err != nil {
		return otto.NullValue()
	}

	v, _ := otto.ToValue(ips[0].String())
	return v
}

func myIPAddress(call otto.FunctionCall) otto.Value {
	ifs, err := net.Interfaces()
	if err != nil {
		return otto.NullValue()
	}

	for _, ifn := range ifs {
		if ifn.Flags&net.FlagUp != net.FlagUp {
			continue
		}

		addrs, err := ifn.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if ok && ip.IP.IsGlobalUnicast() {
				ipstr := ip.IP.String()
				v, _ := otto.ToValue(ipstr)
				return v
			}
		}
	}
	return otto.NullValue()
}
