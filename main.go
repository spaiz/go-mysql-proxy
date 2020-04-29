package main

import (
	"context"
	proxy2 "go-mysql-proxy/proxy"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	proxy := proxy2.NewProxy("127.0.0.1", ":3306", ctx)
	proxy.EnableDecoding()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			log.Printf("Signal received %v, stopping and exiting...", sig)
			cancel()
		}
	}()

	err := proxy.Start("3336")
	if err != nil {
		log.Fatal(err)
	}
}
