# JA4H Fingerprint in Go

An implementation of the [JA4H hash algorithm](https://github.com/FoxIO-LLC/ja4) in Go.

## JA4H_b Issues

Note that this is not a perfect implementation of the algorithm. 

The JA4H_b section will not be correct because the fingerprint should be the truncated SHA256 hash of the request headers in the order they appear.

Since Go stores the headers in a map, it does not keep the ordering as they appeared in the request.

This implementation of the JA4H_b section sorts the request headers before hashing to make the fingerprint consistent.

## Example Middleware

Below is an example of using the library as a middleware. It calculates the JA4H fingerprint and adds it as a new header.

```go
package main

import (
	"fmt"
	"net/http"

	ja4h "github.com/lum8rjack/go-ja4h"
)

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index page with a new JA4H header: %s\n", r.Header.Get("JA4H"))
}

func ja4hMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("JA4H", ja4h.JA4H(r))
		nextHandler.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", ja4hMiddleware(http.HandlerFunc(indexPageHandler)))
	http.ListenAndServe(":8080", mux)
}
```

After starting the server, you can perform different requests to view the fingerprint.

```bash
$ curl http://127.0.0.1:8080
Index page with a new JA4H header: ge11nn020000_a00508f53a24_000000000000_000000000000

$ curl http://127.0.0.1:8080 -H "Accept-Language: en-us"
Index page with a new JA4H header: ge11nn03enus_7cf2b917f4b0_000000000000_000000000000

$ curl http://127.0.0.1:8080 -H "Accept-Language: en-us" -H "Cookie: admin=true" -X POST
Index page with a new JA4H header: po11cn03enus_55e041b6e2b4_8c6976e5b541_8c6976e5b541
```

# References

- [JA4+ Network Fingerprinting](https://blog.foxio.io/ja4+-network-fingerprinting)
- [ja4 GitHub](https://github.com/FoxIO-LLC/ja4)
