package main

import (
	"context"
	"fmt"
	"net"
	"strings"
)

// GET /index.html HTTP/1.1
// Host: localhost:4221
// User-Agent: curl/7.64.1

type httpRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

// HTTP/2 200
// content-type: text/html; charset=UTF-8
// expires: Sat, 18 May 2024 11:30:33 GMT
// content-length: 24
// <!doctype html>
// <html>
// <head>
//
// <title>Example Domain</title>
type httpResponse struct {
	Version    string
	StatusCode string
	Reason     string
	Headers    map[string]string
	Body       string
}

func handleConnection(cxt context.Context, conn net.Conn) {

	defer conn.Close()

	http_request := httpRequest{}
	err := parse_request(conn, &http_request)
	if err != nil {
		fmt.Println("error parsing request; ", err)
	}

	response := httpResponse{}
	prepareResponse(cxt, http_request, &response)

	err = send_response(conn, response)
	if err != nil {
		fmt.Println("error when sending response: ", err)
	}

	fmt.Printf("\nRequest:  %+v\n", http_request)
	fmt.Printf("Response: %+v\n", response)
}

func parse_request(conn net.Conn, request *httpRequest) error {
	BUF_SIZE := 256
	temp_buff := make([]byte, BUF_SIZE)
	request_bytes := make([]byte, 0)
	for {
		n, err := conn.Read(temp_buff)
		fmt.Println("Read ", n, "Bytes")
		if err != nil {
			fmt.Println("Error when reading request;", err)
		}
		if n < BUF_SIZE {
			temp_buff = temp_buff[:n]
			request_bytes = append(request_bytes, temp_buff...)
			break
		}
		request_bytes = append(request_bytes, temp_buff...)

	}
	fmt.Println("Received ", len(request_bytes), "Bytes")

	request_string := string(request_bytes)
	splitted_request := strings.Split(request_string, HTTP_EOL)
	status_line := splitted_request[0]
	splitted_status_line := strings.Split(status_line, " ")
	request.Method, request.Path, request.Version = splitted_status_line[0], splitted_status_line[1], splitted_status_line[2]
	request.Body = splitted_request[len(splitted_request)-1]
	request.Body = strings.TrimSpace(request.Body)

	http_headers := splitted_request[1 : len(splitted_request)-2]
	http_headers_map := make(map[string]string)

	for i := range http_headers {
		splitted_header := strings.SplitN(http_headers[i], ":", 2)
		if len(splitted_header) == 2 {
			http_headers_map[strings.ToLower(splitted_header[0])] = strings.Trim(splitted_header[1], " ")
		}
	}
	request.Headers = http_headers_map

	return nil
}

func prepareResponse(cxt context.Context, http_request httpRequest, http_response *httpResponse) {
	http_response.Version = HTTPVersion

	http_response.Headers = make(map[string]string)

	switch {
	case http_request.Path == "/":
		{
			prepareRootResponse(cxt, http_request, http_response)
		}
	case strings.HasPrefix(http_request.Path, "/echo/"):
		{
			prepareEchoResponse(cxt, http_request, http_response)
		}
	case http_request.Path == "/user-agent":
		{
			prepareUserAgentEndpointResponse(cxt, http_request, http_response)
		}
	case strings.HasPrefix(http_request.Path, "/files/") && strings.ToLower(http_request.Method) == "get":
		{
			prepareGetFileResponse(cxt, http_request, http_response)
		}
	case strings.HasPrefix(http_request.Path, "/files/") && strings.ToLower(http_request.Method) == "post":
		{
			preparePostFileResponse(cxt, http_request, http_response)
		}
	default:
		{
			prepareUnknownResponse(cxt, http_request, http_response)
		}

	}
	compressionMiddleWare(cxt, http_request, http_response)

}
func send_response(conn net.Conn, response httpResponse) error {

	status_line := fmt.Sprintf("%s %s %s", HTTPVersion, response.StatusCode, response.Reason) + HTTP_EOL
	header_line := ""
	for k, v := range response.Headers {
		header_line = header_line + k + ": " + v + HTTP_EOL
	}
	header_line += HTTP_EOL
	http_response_string := status_line + header_line + response.Body

	fmt.Println("Writing response: ", http_response_string)
	n, err := conn.Write([]byte(http_response_string))
	fmt.Println("Sent ", n, "Bytes")
	if err != nil {
		fmt.Println("Error writing response", err)
		return err
	}
	return nil
}

func compressionMiddleWare(cxt context.Context, http_request httpRequest, http_response *httpResponse) {

	if _, ok := http_request.Headers["accept-encoding"]; !ok {
		return
	}
	var compression string = http_request.Headers["accept-encoding"]
	switch compression {
	case "gzip":
		{
			http_response.Headers["Content-Encoding"] = "gzip"
		}
	default:
		{
			return
		}
	}
}
