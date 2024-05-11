package main

import (
	"fmt"
	"strings"
)

func prepareEchoResponse(http_request httpRequest, http_response *httpResponse) {
	http_response.StatusCode = "200"
	http_response.Reason = "OK"

	splitted_path := strings.Split(http_request.Path, "/")
	http_response.Body = splitted_path[len(splitted_path)-1]
	http_response.Headers["Content-Type"] = "text/plain"
	http_response.Headers["Content-Length"] = fmt.Sprintf("%d", len(http_response.Body))

}

func prepareUserAgentEndpoint(http_request httpRequest, http_response *httpResponse) {

	if _, ok := http_request.Headers["user-agent"]; !ok {
		http_response.StatusCode = "400"
		http_response.Reason = "Bad Request"
		return
	}

	http_response.StatusCode = "200"
	http_response.Reason = "OK"

	http_response.Body = http_request.Headers["user-agent"]
	http_response.Headers["Content-Type"] = "text/plain"
	http_response.Headers["Content-Length"] = fmt.Sprintf("%d", len(http_response.Body))
}

func prepareUnknownResponse(http_request httpRequest, http_response *httpResponse) {
	http_response.StatusCode = "404"
	http_response.Reason = "Not Found"
	http_response.Version = HTTPVersion
}

func prepareRootResponse(http_request httpRequest, http_response *httpResponse) {
	http_response.StatusCode = "200"
	http_response.Reason = "OK"
}
