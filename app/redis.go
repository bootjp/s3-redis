package app

import (
	"context"
	"strings"

	"github.com/tidwall/redcon"
)

type Redis struct {
	port          string
	host          string
	UseCache      bool
	UseDirtyWrite bool
	s3            *S3
	handler       func(conn redcon.Conn, cmd redcon.Command)
}

func NewRedis(s3 *S3, config Config) *Redis {
	return &Redis{
		port:          config.Redis.Port,
		host:          config.Redis.Bind,
		UseCache:      config.UseCache,
		UseDirtyWrite: config.UseDirtyWrite,
		s3:            s3,
	}
}

func (r *Redis) Set(ctx context.Context, key string, value []byte) error {
	return r.s3.Put(ctx, key, value)
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	return r.s3.Get(ctx, key)
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.s3.Delete(ctx, key)
}

func (r *Redis) Keys(ctx context.Context, key string) ([]string, error) {
	pattern := strings.HasSuffix(key, "*")
	key = strings.TrimSuffix(key, "*")

	keys, err := r.s3.List(ctx, key)
	if err != nil {
		return nil, err
	}

	var result []string

	// emulate redis keys
	for _, s := range keys {
		if !pattern {
			if key != s {
				continue
			}
			result = append(result, s)
			continue
		}
		result = append(result, s)
	}

	return result, nil
}

func (r *Redis) Handler(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
		return
	case "info":
		conn.WriteBulkString(`
# Server
s3-redis version: 0.0.0
`)
		return
	case "command":
		conn.WriteString("OK")
		return
	case "ping":
		conn.WriteString("PONG")
		return
	case "set":
		if len(cmd.Args) != 3 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		err := r.Set(context.Background(), string(cmd.Args[1]), cmd.Args[2])
		if err != nil {
			conn.WriteError(err.Error())
			return
		}
		conn.WriteString("OK")
		return

	case "get":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		value, err := r.Get(context.Background(), string(cmd.Args[1]))
		if err != nil {
			conn.WriteError(err.Error())
			return
		}
		conn.WriteBulk(value)
		return

	case "del":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		err := r.Delete(context.Background(), string(cmd.Args[1]))
		if err != nil {
			conn.WriteError(err.Error())
			return
		}
		conn.WriteString("OK")
		return
	case "keys":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		keys, err := r.Keys(context.Background(), string(cmd.Args[1]))
		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		conn.WriteArray(len(keys))
		for _, key := range keys {
			conn.WriteBulkString(key)
		}
		return
	case "exists":
		keys, err := r.Get(context.Background(), string(cmd.Args[1]))
		if err != nil {
			conn.WriteInt(0)
			return
		}

		var res int

		if len(keys) > 0 {
			res = 1
		} else {
			res = 0
		}
		conn.WriteInt(res)
		return
	}

}
