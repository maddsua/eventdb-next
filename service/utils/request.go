package utils

import (
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const defaultXffHeader = "X-Forwarded-For"
const defaultXridHeader = "X-Request-Id"
const reqIddByteLength = 4

func RequestID(writer http.ResponseWriter, req *http.Request) string {

	if rid := strings.TrimSpace(req.Header.Get(defaultXridHeader)); rid != "" {
		return rid
	}

	buff := make([]byte, reqIddByteLength)
	rand.Read(buff)
	rid := fmt.Sprintf("%x-MWS", buff)

	req.Header.Set(defaultXridHeader, rid)
	writer.Header().Set(defaultXridHeader, rid)

	return rid
}

func ReadRequestID(req *http.Request) string {
	return req.Header.Get(defaultXridHeader)
}

func ClientIP(req *http.Request) string {

	if ipString := parseXFFHeader(req.Header); len(ipString) > 0 {
		return ipString
	}

	requestHost, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return ""
	}

	return requestHost
}

func parseXFFHeader(headers http.Header) string {

	xffHeader := headers.Get(defaultXffHeader)
	if len(xffHeader) == 0 {
		return ""
	}

	lastComma := strings.LastIndex(xffHeader, ",")
	if lastComma < 0 {
		return xffHeader
	}

	lastHop := strings.TrimSpace(xffHeader[lastComma+1:])
	if len(lastHop) == 0 {
		return ""
	}

	return lastHop
}
