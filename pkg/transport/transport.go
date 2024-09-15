package transport

import "net"

type Transport interface {
	ListenAndAccept() error

	Consume() <-chan net.Conn

	CLose() error
}
