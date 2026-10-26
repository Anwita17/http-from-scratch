package http1

import (
	"fmt"
	"net"
	"net/http"
)

type Response struct {
	conn        net.Conn
	header      http.Header
	statusCode  int
	wroteHeader bool
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
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.conn.Write(data)
}

func (r *Response) WriteHeader(statusCode int) {
	//panic("Response.WriteHeader: Not Implemented")
	if r.wroteHeader {
		return
	}
	r.wroteHeader = true
	r.statusCode = statusCode

	statusText := http.StatusText(statusCode)
	if statusText == "" {
		statusText = "status"
	}
	fmt.Fprintf(r.conn, "HTTP/1.1 %d %s\r\n", statusCode, statusText)

	if r.header.Get("Connection") == "" {
		r.header.Set("Connection", "close")
	}

	r.header.Write(r.conn)

	fmt.Fprintf(r.conn, "\r\n")
}
