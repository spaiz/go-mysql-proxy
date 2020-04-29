package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
)

func NewProxy(host, port string, ctx context.Context) *Proxy {
	return &Proxy{
		host: host,
		port: port,
		ctx: ctx,
	}
}

type Proxy struct {
	host           string
	port           string
	connectionId   uint64
	enableDecoding bool
	ctx context.Context
	shutDownAsked bool
}

func (r *Proxy) Start(port string) error {
	log.Printf("Start listening on: %s", port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	go func() {
		log.Printf("Waiting for shut down signal ^C")
		<-r.ctx.Done()
		r.shutDownAsked = true
		log.Printf("Shut down signal received, closing connections...")
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		r.connectionId += 1
		if err != nil {
			log.Printf("Failed to accept new connection: [%d] %s", r.connectionId, err.Error())
			if r.shutDownAsked {
				log.Printf("Shutdown asked [%d]", r.connectionId,)
				break
			}
			continue
		}

		log.Printf("Connection accepted: [%d] %s", r.connectionId, conn.RemoteAddr())
		go r.handle(conn, r.connectionId, r.enableDecoding)
	}

	return nil
}

func (r *Proxy) handle(conn net.Conn, connectionId uint64, enableDecoding bool) {
	connection := NewConnection(r.host, r.port, conn, connectionId, enableDecoding)
	err := connection.Handle()
	if err != nil {
		log.Printf("Error handling proxy connection: %s", err.Error())
	}
}

func (r *Proxy) EnableDecoding() {
	r.enableDecoding = true
}