package http1

import (
	"bufio"
	"bytes"
	"errors"
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
	line, err := readLine(reader)
	if err != nil {
		return nil, fmt.Errorf("Request line: %w", err)
	}
	method, uri, version, err := ParseRequestLine(line)
	if err != nil {
		return nil, err
	}
	if err := checkProtoVersion(version); err != nil {
		return nil, err
	}
	log.Printf("method=%q uri=%q version=%q", method, uri, version)
	// Header
	header_store := http.Header{}
	for {
		header, err := readLine(reader)
		if err != nil {
			return nil, fmt.Errorf("Error in header: %w", err)
		}
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
	hosts := header_store.Values("Host")
	if version == "HTTP/1.1" && len(hosts) == 0 {
		return nil, fmt.Errorf("missing host header!")
	}
	if len(hosts) > 1 {
		return nil, fmt.Errorf("multiple host headers")
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
	if te := header_store.Get("Transfer-Encoding"); te != "" {
		if !strings.EqualFold(te, "chunked") {
			return nil, fmt.Errorf("unsupported transfer-encoding %q", te)
		}
		body, err := readChunkedBody(reader)
		if err != nil {
			return nil, err
		}
		request.ContentLength = int64(len(body))
		request.Body = io.NopCloser(bytes.NewReader(body))
	} else if cl := header_store.Get("Content-Length"); cl != "" {
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
	if !isToken(header_split[0]) {
		return "", "", "", fmt.Errorf("Request is malformed: %q", line)
	}
	return header_split[0], header_split[1], header_split[2], nil
}

// getting the key:value pairs
func ParseRequestHeader(line string) (name, val string, err error) {
	name, val, seperator := strings.Cut(line, ":")
	if !seperator || !isToken(name) {
		return "", "", fmt.Errorf("Header is malformed: %q", line)
	}
	return name, strings.Trim(val, " \t"), nil
}

// check all protocol version according to RFC
var errWithVersion = errors.New("http version not supported")

func checkProtoVersion(version string) error {
	protocol, num, found := strings.Cut(version, "/")
	if !found || protocol != "HTTP" {
		return fmt.Errorf("malformed version %q", version)
	}

	major, minor, found := strings.Cut(num, ".")
	if !found || !isDigit(major) || !isDigit(minor) {
		return fmt.Errorf("malformed version %q", version)
	}
	if version != "HTTP/1.0" && version != "HTTP/1.1" {
		return fmt.Errorf("%w : %q", errWithVersion, version)
	}
	return nil
}

func isDigit(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func isToken(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case strings.IndexByte("!#$%&'*+-.^_`/~", c) >= 0:
		default:
			return false
		}
	}
	return true
}

// reads a line and requires it to end with CRLF
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(line, "\r\n") {
		return "", fmt.Errorf("line not terminated by CRLF: %q", line)
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}

func readChunkedBody(reader *bufio.Reader) ([]byte, error) {
	var body bytes.Buffer
	for {
		line, err := readLine(reader)
		if err != nil {
			return nil, fmt.Errorf("chunksize line: %w", err)
		}
		sizeStr, _, _ := strings.Cut(line, ";")
		size, err := strconv.ParseUint(sizeStr, 16, 63)
		if err != nil {
			return nil, fmt.Errorf("invalid chunk size %q", sizeStr)
		}
		if size == 0 {
			break
		}
		if _, err := io.CopyN(&body, reader, int64(size)); err != nil {
			return nil, fmt.Errorf("chunk data: %w", err)
		}
		if end, err := readLine(reader); err != nil || end != "" {
			return nil, fmt.Errorf("chunk not followed by CRLF")
		}
	}
	for {
		line, err := readLine(reader)
		if err != nil {
			return nil, fmt.Errorf("trailer: %w", err)
		}
		if line == "" {
			break
		}
	}
	return body.Bytes(), nil
}
