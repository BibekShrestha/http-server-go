package main

import (
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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	http_request := httpRequest{}
	err := parse_request(conn, &http_request)
	if err != nil {
		fmt.Println("error parsing request; ", err)
	}

	response := httpResponse{}
	prepareResponse(http_request, &response)

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
		request_bytes = append(request_bytes, temp_buff...)
		if n < BUF_SIZE {
			break
		}
	}
	fmt.Println("Received ", len(request_bytes), "Bytes")

	request_string := string(request_bytes)
	splitted_request := strings.Split(request_string, HTTP_EOL)
	status_line := splitted_request[0]
	splitted_status_line := strings.Split(status_line, " ")
	request.Method, request.Path, request.Version = splitted_status_line[0], splitted_status_line[1], splitted_status_line[2]
	request.Body = splitted_request[len(splitted_request)-1]
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

func prepareResponse(http_request httpRequest, http_response *httpResponse) {
	http_response.Version = HTTPVersion

	http_response.Headers = make(map[string]string)

	switch {
	case http_request.Path == "/":
		{
			prepareRootResponse(http_request, http_response)
		}
	case strings.HasPrefix(http_request.Path, "/echo/"):
		{
			prepareEchoResponse(http_request, http_response)
		}
	case http_request.Path == "/user-agent":
		{
			prepareUserAgentEndpoint(http_request, http_response)
		}
	default:
		{
			prepareUnknownResponse(http_request, http_response)
		}

	}

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