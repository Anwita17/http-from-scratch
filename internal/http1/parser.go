package http1

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {
	//panic("ParseRequest: Not Implemented")
	//buffer := make([]byte, 1024)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Request line: %w", err)
	}
	line = strings.TrimRight(line, "\r\n")
	method, uri, version, err := ParseRequestLine(line)
	if err != nil {
		return nil, err
	}
	log.Printf("method=%q uri=%q version=%q", method, uri, version)
	// Header
	header_store := http.Header{}
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("Error in header: %w", err)
		}
		header = strings.TrimRight(header, "\r\n")
		if header == "" {
			log.Printf("Reached the end of Headers!")
			break
		}
		name, val, err := ParseRequestHeader(header)
		if err != nil {
			return nil, err
		}
		header_store.Add(name, val)
		log.Printf("name= %q, value=%q", name, val)
	}
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, fmt.Errorf("Bad URI %q: %w", uri, err)
	}
	//	if version != "HTTP/1.1" && version != "HTTP/1.0" {
	//		return nil, fmt.Errorf("Unsupported version of HTTP is used : %q", version)
	//	}
	//	if version == "HTTP/1.1" && header_store.Get("Host") == "" {
	//		return nil, fmt.Errorf("Missing Host in Header!")
	//	}
	request := &http.Request{
		Method:     method,
		URL:        u,
		Proto:      version,
		Header:     header_store,
		Host:       header_store.Get("Host"),
		RequestURI: uri,
		Body:       http.NoBody,
	}

	log.Printf("parsed request : %s %s host=%s headers=%s", request.Method, request.URL.Path, request.Host, request.Header)
	return request, nil
}

func ParseRequestLine(line string) (method, uri, version string, err error) {
	header_split := strings.Split(line, " ")
	if len(header_split) != 3 {
		return "", "", "", fmt.Errorf("Request is malformed: %q", line)
	}
	return header_split[0], header_split[1], header_split[2], nil
}

// getting the key:value pairs
func ParseRequestHeader(line string) (name, val string, err error) {
	name, val, seperator := strings.Cut(line, ":")
	if !seperator || name == "" || strings.ContainsAny(name, " \t") {
		return "", "", fmt.Errorf("Header is malformed: %q", line)
	}
	return name, strings.Trim(val, " \t"), nil
}
