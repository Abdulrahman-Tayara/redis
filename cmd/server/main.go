package main

import (
	"redis/internal/server"
	"redis/internal/transport"
)

func main() {

	s := server.NewRedisServer(transport.NewTcpTransport(":9871"))

	if err := s.Serve(); err != nil {
		panic(err)
	}

}
