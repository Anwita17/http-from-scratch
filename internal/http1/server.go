package http1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func ListenAndServe(addr string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer listener.Close()

	if handler == nil {
		handler = http.DefaultServeMux
	}

	for {
		//panic("ListenAndServe: Not implemented")
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("Error in Connection : %w", err)
		}
		go func() {
			if err := HandleConnection(conn, handler); err != nil {
				log.Printf("Connection error %s: %v", conn.RemoteAddr(), err)
			}
		}()

	}
}

func HandleConnection(conn net.Conn, mux http.Handler) error {
	//panic("HandleConnection: Not Implemented")
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			_ = fmt.Errorf("Close connection: %w", err)
		}
	}(conn)

	buffer := bufio.NewReader(conn)
	//n, err := conn.Read(buffer)
	for {
		if _, err := buffer.Peek(1); err != nil {
			return nil
		}
		request, err := ParseRequest(buffer)
		if err != nil {
			return err
		}

		closeAfter := shouldClose(request)
		w := NewResponse(conn)
		if closeAfter {
			w.Header().Set("Connection", "close")
		}
		mux.ServeHTTP(w, request)
		if err := w.finish(request); err != nil {
			return err
		}

		if closeAfter {
			return nil
		}
		//return nil
	}
	//_, err = conn.Write([]byte("HTTP/1.1 200 OK \r\n Content-Length: 0\r\n"))
}

func shouldClose(request *http.Request) bool {
	conn := request.Header.Get("Connection")
	if request.Proto == "HTTP/1.0" {
		return !strings.EqualFold(conn, "keep-alive")
	}
	return strings.EqualFold(conn, "close")
}
