package main

import "log"

func main() {
	proxy := NewProxy("127.0.0.1", ":3306")
	err := proxy.Start("3336")
	if err != nil {
		log.Fatal(err)
	}
}
