package session

import (
	"io"
	"net"
	"net/http"
	"strings"
)

func copyResponse(w http.ResponseWriter, res *http.Response) error {
	// copy header
	for key, values := range w.Header() {
		for _, value := range values {
			res.Header.Add(key, value)
		}
	}

	w.WriteHeader(res.StatusCode)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	_, err := io.Copy(w, res.Body)
	return err
}

type CustomCopyHeader func(r *http.Request, output http.Header)

func getForwardingHeader(req *http.Request, copyHeader CustomCopyHeader) http.Header {
	reqHeader := req.Header
	result := http.Header{}

	for _, cookie := range reqHeader[http.CanonicalHeaderKey("Cookie")] {
		result.Add("Cookie", cookie)
	}
	if req.Host != "" {
		result.Set("Host", req.Host)
	}
	if origin := reqHeader.Get("Origin"); origin != "" {
		result.Add("Origin", origin)
	}
	for _, protocol := range reqHeader[http.CanonicalHeaderKey("Sec-Websocket-Protoco")] {
		result.Add("Sec-Websocket-Protocol", protocol)
	}
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// preserve previous X-Forwarded-For
		if oldClientIPs, ok := req.Header["X-Forwarded-For"]; ok {
			result.Set("X-Forwarded-For",
				strings.Join(append(oldClientIPs, clientIP), ", "))
		} else {
			result.Set("X-Forwarded-For", clientIP)
		}
	}
	if req.TLS != nil {
		result.Set("X-Forwarded-Proto", "https")
	} else {
		result.Set("X-Forwarded-Proto", "http")
	}

	if copyHeader != nil {
		copyHeader(req, result)
	}
	return result
}

func copyHeaders(srcHeader http.Header, destHeader *http.Header, keys []string) {
	for _, key := range keys {
		if value := srcHeader.Get(key); value != "" {
			destHeader.Set(key, value)
		}
	}
}
