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

func prepareGetFileResponse(cxt context.Context, http_request httpRequest, http_response *httpResponse) {
	splitted_path := strings.SplitN(http_request.Path, "/", 3)
	if len(splitted_path) < 3 {
		fmt.Printf("Invalid input for file response; %+v -> %+v\n", http_request.Path, splitted_path)
		prepareUnknownResponse(cxt, http_request, http_response)
		return
	}

	fileName := fmt.Sprintf("%v", cxt.Value("workDir")) + string(os.PathSeparator) + splitted_path[2]
	fileContent, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Printf("Error when opening file: %s: %+v\n", fileName, err)
		prepareUnknownResponse(cxt, http_request, http_response)
		return
	}
	http_response.StatusCode = "200"
	http_response.Reason = "OK"
	http_response.Body = string(fileContent)
	http_response.Headers["Content-Type"] = "application/octet-stream"
	http_response.Headers["Content-Length"] = fmt.Sprintf("%d", len(http_response.Body))

}

func preparePostFileResponse(cxt context.Context, http_request httpRequest, http_response *httpResponse) {
	splitted_path := strings.SplitN(http_request.Path, "/", 3)
	if len(splitted_path) < 3 {
		fmt.Printf("Invalid input for file response; %+v -> %+v\n", http_request.Path, splitted_path)
		prepareUnknownResponse(cxt, http_request, http_response)
		return
	}

	fileName := fmt.Sprintf("%v", cxt.Value("workDir")) + string(os.PathSeparator) + splitted_path[2]
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Printf("Error when opening file for writing: %s: %+v\n", fileName, err)
		prepareUnknownResponse(cxt, http_request, http_response)
		return
	}
	if n, err := file.WriteString(http_request.Body); n != len(http_request.Body) || err != nil {
		fmt.Println("Error when writing to file", err)
		prepareUnknownResponse(cxt, http_request, http_response)
		return
	}

	http_response.StatusCode = "201"
	http_response.Reason = "Created"

}
