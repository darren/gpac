## gpac

[![GoDoc](https://godoc.org/github.com/darren/gpac?status.png)](https://godoc.org/github.com/darren/gpac)


This package provides a pure Go [pac](https://developer.mozilla.org/en-US/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_(PAC)_file) parser based on [otto](https://github.com/robertkrimen/otto)

## Example usage

```go
 pac, _ := gpac.New(`
     function FindProxyForURL(url, host) {
         if (isPlainHostName(host)) return DIRECT;
         else return "PROXY 127.0.0.1:8080";
     }
 `)
 r,_ := pac.FindProxyForURL("http://www.example.com/")
 fmt.Println(r)
```
