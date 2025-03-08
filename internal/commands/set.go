package commands

import (
	"errors"
	"redis/internal/server"
	"redis/pkg/argparser"
	"redis/pkg/ds"
	"redis/pkg/iox"
	"redis/pkg/utils"
	"time"
)

var setCommandOptionsSchema = []argparser.ArgInfo{
	{"NX", true}, {"XX", true},
	{"GET", true},
	{"EX", false}, {"PX", false},
	{"PXAT", false},
	{"KEEPTTL", true},
}

func (s *Server) HandleSet() server.CommandHandlerFunc {

	return func(ctx *server.Context, w iox.AnyWriter) {
		args := ctx.Args()

		if len(args) < 2 {
			_, _ = w.WriteAny(errors.New("ERR wrong number of arguments for 'set' command"))
			return
		}

		key, value := args[0], args[1]

		opts, err := argparser.Parse(args[2:], setCommandOptionsSchema)
		if err != nil {
			_, _ = w.WriteAny(err)
			return
		}

		setOptions := &ds.SetOptions{
			KeepTTL:        opts.GetOrDefault("KEEPTTL", false).(bool),
			SetIfNotExists: opts.GetOrDefault("NX", false).(bool),
			SetIfExists:    opts.GetOrDefault("XX", false).(bool),
			ExpireAt: calculateTTL(
				utils.MustParseInt(opts.GetOrDefault("PX", 0)),
				utils.MustParseInt(opts.GetOrDefault("EX", 0)),
				utils.MustParseInt(opts.GetOrDefault("PXAT", 0)),
				utils.MustParseInt(opts.GetOrDefault("EXAT", 0)),
			),
		}

		hashTable := s.store.HashTable()

		if err = hashTable.Set(key.(string), value, setOptions); err != nil {
			_, _ = w.WriteAny(err)
			return
		}

		if get := opts.GetOrDefault("GET", false).(bool); get {
			if val, ok := hashTable.Get(key.(string)); !ok {
				_, _ = w.WriteAny(nil)
				return
			} else {
				_, _ = w.WriteAny(val)
			}
		} else {
			_, _ = w.WriteAny("OK")
		}
	}
}

func calculateTTL(msDuration int, secDuration int, msUnixTime int, secUnixTime int) int64 {
	if msDuration > 0 {
		return time.Now().Add(time.Millisecond * time.Duration(msDuration)).UnixMilli()
	}
	if secDuration > 0 {
		return time.Now().Add(time.Second * time.Duration(secDuration)).UnixMilli()
	}
	if msUnixTime > 0 {
		return int64(msUnixTime)
	}
	if secUnixTime > 0 {
		return int64(secUnixTime * 1000)
	}
	return 0
}
