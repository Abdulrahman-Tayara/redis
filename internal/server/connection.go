package server

import (
	"net"
	"redis/pkg/resp"
)

type RedisConnection struct {
	net.Conn
	resp.CommandReader
	*resp.Writer

	Id int32
}

func NewRedisConnection(conn net.Conn, id int32) *RedisConnection {
	return &RedisConnection{
		CommandReader: resp.NewCommandReader(conn),
		Writer:        resp.NewRespWriter(conn),
		Conn:          conn,
		Id:            id,
	}
}
