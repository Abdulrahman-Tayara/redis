package commands

import (
	"errors"
	"redis/internal/server"
	"redis/pkg/iox"
)

func (s *Server) HandleSet() server.CommandHandlerFunc {
	return func(ctx *server.Context, w iox.AnyWriter) {
		args := ctx.Args()

		if len(args) != 2 {
			_, _ = w.WriteAny(errors.New("ERR wrong number of arguments for 'set' command"))
			return
		}

		key, value := args[0], args[1]

		hashTable := s.store.HashTable()

		if err := hashTable.Set(key.(string), value); err != nil {
			_, _ = w.WriteAny(err)
			return
		}

		_, _ = w.WriteAny("OK")
	}
}
