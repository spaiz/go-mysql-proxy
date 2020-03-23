package main

import (
	"fmt"
	"log"
	"net"
)

func NewProxy(host, port string) *Proxy {
	return &Proxy{
		host: host,
		port: port,
	}
}

type Proxy struct {
	host string
	port string
	connectionId uint64
}

func (r *Proxy) Start(port string) error {
	log.Printf("Start listening on: %s", port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		r.connectionId += 1
		log.Printf("Connection accepted: [%d] %s", r.connectionId, conn.RemoteAddr())
		if err != nil {
			log.Printf("Failed to accept new connection: [%d] %s", r.connectionId, err.Error())
			continue
		}

		go r.handle(conn, r.connectionId)
	}
}

func (r *Proxy) handle(conn net.Conn, connectionId uint64) {
	connection := NewConnection(r.host, r.port, conn, connectionId)
	err := connection.Handle()
	if err != nil {
		log.Printf("Error handling proxy connection: %s", err.Error())
	}
}