// Package ja4h is a Go implementation of the JA4HTTP (JA4H) hashing algorithm
// (https://github.com/FoxIO-LLC/ja4).
//
// Note:
//
// This is not a perfect implementation of the algorithm. The JA4H_b section will not be correct
// because the fingerprint should be the truncated SHA256 hash of the request headers in the
// order they appear.
//
// Since Go stores the request headers in a map, it does not keep the ordering as they appeared
// in the request. This implementation of the JA4H_b section sorts the headers before hashing to
// make the fingerprint consistent.
package ja4h

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func http_method(method string) string {
	return strings.ToLower(method)[:2]
}

func http_version(version string) string {
	v := strings.Split(version, "/")
	if len(v) == 2 {
		if v[1] == "2" || v[1] == "2.0" {
			return "20"
		}
	}

	return "11"
}

func hasCookie(req *http.Request) string {
	if len(req.Cookies()) > 0 {
		return "c"
	}
	return "n"
}

func hasReferer(referer string) string {
	if referer != "" {
		return "r"
	}
	return "n"
}

func num_headers(headers http.Header) int {
	len_headers := len(headers)
	if headers.Get("Cookie") != "" {
		len_headers--
	}
	if headers.Get("Referer") != "" {
		len_headers--
	}
	return len_headers
}

func language(headers http.Header) string {
	lan := headers.Get("Accept-Language")
	if lan != "" {
		clean := strings.ReplaceAll(lan, "-", "")
		lower := strings.ToLower(clean)
		first := strings.Split(lower, ",")[0] + "0000"
		return first[:4]
	}
	return "0000"
}

// 1. HTTP Method, GET="ge", PUT="pu", POST="po", etc.
//
// 2. HTTP Version, 2.0="20", 1.1="11"
//
// 3. Cookie, if there's a Cookie "c", if no Cookie "n"
//
// 4. Referer, if there's a Referer "r", if no Referer "n"
//
// 5. Number of HTTP Headers (ignore Cookie and Referer)
//
// 6. First 4 characters of primary Accept-Language (0000 if no Accept-Language)
func JA4H_a(req *http.Request) string {
	method := http_method(req.Method)
	version := http_version(req.Proto)
	cookie := hasCookie(req)
	referer := hasReferer(req.Referer())
	num_headers := num_headers(req.Header)
	accept_lang := language(req.Header)

	return fmt.Sprintf("%s%s%s%s%02d%s", method, version, cookie, referer, num_headers, accept_lang)
}

// Truncated SHA256 hash of Headers, in the order they appear
//
// ISSUE: Go HTTP request headers are a map and does not keep the ordering
func JA4H_b(req *http.Request) string {
	ordered_headers := make([]string, 0, len(req.Header))
	for h := range req.Header {
		ordered_headers = append(ordered_headers, h)
	}
	sort.Strings(ordered_headers)
	allheaders := strings.Join(ordered_headers, "")

	hash := sha256.New()
	hash.Write([]byte(allheaders))
	bs := hash.Sum(nil)
	return fmt.Sprintf("%x", bs)[:12]
}

// Truncated SHA256 hash of Cookie Fields, sorted
func JA4H_c(req *http.Request) string {
	if len(req.Cookies()) == 0 {
		return strings.Repeat("0", 12)
	}
	ordered_cookies := make([]string, 0, len(req.Cookies()))
	for _, c := range req.Cookies() {
		ordered_cookies = append(ordered_cookies, c.Name)
	}
	sort.Strings(ordered_cookies)
	allcookies := strings.Join(ordered_cookies, "")

	hash := sha256.New()
	hash.Write([]byte(allcookies))
	bs := hash.Sum(nil)
	return fmt.Sprintf("%x", bs)[:12]
}

// Truncated SHA256 hash of Cookie Fields + Values, sorted
func JA4H_d(req *http.Request) string {
	if len(req.Cookies()) == 0 {
		return strings.Repeat("0", 12)
	}
	ordered_cookies := make([]string, 0, len(req.Cookies()))
	for _, c := range req.Cookies() {
		ordered_cookies = append(ordered_cookies, c.Name)
	}
	sort.Strings(ordered_cookies)
	allcookies := strings.Join(ordered_cookies, "")

	hash := sha256.New()
	hash.Write([]byte(allcookies))
	bs := hash.Sum(nil)
	return fmt.Sprintf("%x", bs)[:12]
}

// HTTP client fingerprint based on each HTTP request.
func JA4H(req *http.Request) string {
	JA4H_a := JA4H_a(req)
	JA4H_b := JA4H_b(req)
	JA4H_c := JA4H_c(req)
	JA4H_d := JA4H_d(req)

	return fmt.Sprintf("%s_%s_%s_%s", JA4H_a, JA4H_b, JA4H_c, JA4H_d)
}
