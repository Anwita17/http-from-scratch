package http1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
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
	/*
			Functions that might be useful:
			* ParseRequest()
		  * mux.ServeHTTP()
	*/
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			_ = fmt.Errorf("Close connection: %w", err)
		}
	}(conn)

	buffer := bufio.NewReader(conn)
	//n, err := conn.Read(buffer)
	request, err := ParseRequest(buffer)
	if err != nil {
		return err
	}

	w := NewResponse(conn)
	mux.ServeHTTP(w, request)
	w.WriteHeader(http.StatusOK)
	//_, err = conn.Write([]byte("HTTP/1.1 200 OK \r\n Content-Length: 0\r\n"))
	return nil
}
