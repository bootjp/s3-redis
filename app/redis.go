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
	r := &Redis{
		port:          config.Redis.Port,
		host:          config.Redis.Bind,
		UseCache:      config.UseCache,
		UseDirtyWrite: config.UseDirtyWrite,
		s3:            s3,
	}

	r.handler = func(conn redcon.Conn, cmd redcon.Command) {
		r.Handler(conn, cmd)
	}

	return r
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
		} else {
			conn.WriteString("OK")
			return
		}

	case "get":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		value, err := r.Get(context.Background(), string(cmd.Args[1]))
		if err != nil {
			conn.WriteError(err.Error())
			return
		} else {
			conn.WriteBulk(value)
			return
		}

	case "del":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		err := r.Delete(context.Background(), string(cmd.Args[0]))
		if err != nil {
			conn.WriteError(err.Error())
			return
		} else {
			conn.WriteString("OK")
			return
		}
	}
}
