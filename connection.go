package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func NewConnection(host string, port string, conn net.Conn, id uint64) *Connection {
	return &Connection{
		host: host,
		port: port,
		conn: conn,
		id: id,
	}
}

type Connection struct {
	id   uint64
	conn net.Conn
	host string
	port string
}

func (r *Connection) Handle() error {
	address := fmt.Sprintf("%s%s", r.host, r.port)
	mysql, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Failed to connection to MySQL: [%d] %s", r.id, err.Error())
		return err
	}

	go func() {
		copied, err := io.Copy(mysql, r.conn)
		if err != nil {
			log.Printf("Conection error: [%d] %s", r.id, err.Error())
		}

		log.Printf("Connection closed. Bytes copied: [%d] %d", r.id, copied)
	}()

	copied, err := io.Copy(r.conn, mysql)
	if err != nil {
		log.Printf("Connection error: [%d] %s", r.id, err.Error())
		return err
	}

	log.Printf("Connection closed. Bytes copied: [%d] %d", r.id, copied)
	return nil
}