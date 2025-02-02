package server

import (
	"net"
	"redis/pkg/resp"
)

type RedisConnection struct {
	net.Conn
	resp.CommandReader
	*resp.Writer

	info map[string]string

	Id int32
}

func NewRedisConnection(conn net.Conn, id int32) *RedisConnection {
	return &RedisConnection{
		CommandReader: resp.NewCommandReader(conn),
		Writer:        resp.NewRespWriter(conn),
		Conn:          conn,
		Id:            id,
		info:          make(map[string]string),
	}
}

func (c *RedisConnection) SetInfo(key, value string) {
	c.info[key] = value
}

func (c *RedisConnection) GetInfo(key string) (string, bool) {
	v, ok := c.info[key]
	return v, ok
}
