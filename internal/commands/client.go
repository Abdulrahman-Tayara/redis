package commands

import (
	"errors"
	"redis/internal/server"
	"redis/pkg/iox"
	"redis/pkg/utils"
	"strings"
)

func (s *Server) HandleClientCommands() server.CommandHandlerFunc {
	return func(ctx *server.Context, w iox.AnyWriter) {
		args := ctx.Args()
		if len(args) < 1 {
			_, _ = w.WriteAny(errors.New("ERR wrong number of arguments for 'client' command"))
			return
		}
		attribute := args[0].(string)
		if handler, ok := clientHandlers[attribute]; !ok {
			_, _ = w.WriteAny(errors.New("ERR unknown client attribute or subcommand"))
		} else {
			if err := handler(ctx, args[1:]); err != nil {
				_, _ = w.WriteAny(err)
				return
			}
		}
	}
}

var clientHandlers = map[string]func(ctx *server.Context, args []any) error{
	"setinfo": handleSetInfo,
}

func handleSetInfo(ctx *server.Context, args []any) error {
	if len(args) < 2 {
		return errors.New("ERR wrong number of arguments for 'setinfo' command")
	}
	key := args[0].(string)
	value := strings.Join(utils.Map(args[1:], func(a any, i int) string {
		return a.(string)
	}), " ")

	ctx.Set(key, value)
	return nil
}
