package http1

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

type Response struct {
	conn       net.Conn
	header     http.Header
	statusCode int
	body       bytes.Buffer
}

func NewResponse(conn net.Conn) *Response {
	//panic("NewResponse: Not Implemented")
	return &Response{conn: conn, header: http.Header{}}
}

func (r *Response) Header() http.Header {
	//	panic("Response.Header: Not Implemented")
	return r.header
}

func (r *Response) Write(data []byte) (n int, err error) {
	//panic("Response.Write: Not Implemented")
	r.WriteHeader(http.StatusOK)
	return r.body.Write(data)
}

func (r *Response) WriteHeader(statusCode int) {
	//panic("Response.WriteHeader: Not Implemented")
	if r.statusCode == 0 {
		r.statusCode = statusCode
	}
}

func (r *Response) finish(request *http.Request) error {
	r.WriteHeader(http.StatusOK)
	r.header.Set("Content-Length", strconv.Itoa(r.body.Len()))

	fmt.Fprintf(r.conn, "HTTP/1.1 %d %s\r\n", r.statusCode, http.StatusText(r.statusCode))
	r.header.Write(r.conn)
	fmt.Fprint(r.conn, "\r\n")
	if request.Method == http.MethodHead {
		return nil
	}
	_, err := r.conn.Write(r.body.Bytes())
	return err
}
