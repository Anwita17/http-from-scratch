package http1

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	request := &http.Request{
		Method:     method,
		URL:        u,
		Proto:      version,
		Header:     header_store,
		Host:       header_store.Get("Host"),
		RequestURI: uri,
		Body:       http.NoBody,
	}
	// handle post with body
	if cl := header_store.Get("Content-Length"); cl != "" {
		n, err := strconv.ParseInt(cl, 10, 64)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("Bad Content-Length %q", cl)
		}
		body := make([]byte, n)
		if _, err := io.ReadFull(reader, body); err != nil {
			return nil, fmt.Errorf("Error in body : %w", err)
		}
		request.ContentLength = n
		request.Body = io.NopCloser(bytes.NewReader(body))
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
