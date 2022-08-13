package app

import (
	"github.com/tidwall/redcon"
)

type Server struct {
	Config Config
	S3     *S3
	Redis  *Redis
}

func NewServer(config Config, s3 *S3, redis *Redis) *Server {
	return &Server{
		Config: config,
		S3:     s3,
		Redis:  redis,
	}
}

func (s *Server) Run() error {
	laddr := s.Config.Redis.Bind + ":" + s.Config.Redis.Port
	err := redcon.ListenAndServe(laddr, s.Redis.Handler, func(conn redcon.Conn) bool {
		// Use this function to accept or deny the connection.
		// log.Printf("accept: %s", conn.RemoteAddr())
		return true
	}, func(conn redcon.Conn, err error) {
		// This is called when the connection has been closed
		// log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
	})

	return err
}
