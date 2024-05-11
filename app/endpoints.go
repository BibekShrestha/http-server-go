package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func prepareEchoResponse(cxt context.Context, http_request httpRequest, http_response *httpResponse) {
	http_response.StatusCode = "200"
	http_response.Reason = "OK"

	splitted_path := strings.Split(http_request.Path, "/")
	http_response.Body = splitted_path[len(splitted_path)-1]
	http_response.Headers["Content-Type"] = "text/plain"
	http_response.Headers["Content-Length"] = fmt.Sprintf("%d", len(http_response.Body))
}

func prepareUserAgentEndpointResponse(cxt context.Context, http_request httpRequest, http_response *httpResponse) {

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

func prepareUnknownResponse(cxt context.Context, http_request httpRequest, http_response *httpResponse) {
	http_response.StatusCode = "404"
	http_response.Reason = "Not Found"
}

func prepareRootResponse(cxt context.Context, http_request httpRequest, http_response *httpResponse) {
	http_response.StatusCode = "200"
	http_response.Reason = "OK"
}

func prepareFileResponse(cxt context.Context, http_request httpRequest, http_response *httpResponse) {
	splitted_path := strings.SplitN(http_request.Path, "/", 3)
	if len(splitted_path) < 3 {
		fmt.Printf("Invalid input for file response; %+v -> %+v", http_request.Path, splitted_path)
		prepareUnknownResponse(cxt, http_request, http_response)
		return
	}

	fileName := fmt.Sprintf("%v", cxt.Value("workDir")) + string(os.PathSeparator) + splitted_path[2]
	fileContent, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Printf("Error when opening file: %s: %+v", fileName, err)
		prepareUnknownResponse(cxt, http_request, http_response)
		return
	}
	http_response.StatusCode = "200"
	http_response.Reason = "OK"
	http_response.Body = string(fileContent)
	http_response.Headers["Content-Type"] = "application/octet-stream"
	http_response.Headers["Content-Length"] = fmt.Sprintf("%d", len(http_response.Body))

}
