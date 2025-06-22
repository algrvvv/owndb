package database

import (
	"bufio"
	"fmt"
	"net"

	"github.com/algrvvv/owndb/internal/config"
	"github.com/algrvvv/owndb/internal/dsl"
	"github.com/algrvvv/owndb/internal/exec"
	"github.com/algrvvv/owndb/internal/storage"
	"github.com/algrvvv/owndb/internal/wal"
	"github.com/rs/zerolog"
)

type DBServer struct {
	storage storage.Storage
	inline  exec.Executor
	port    string
	config  *config.Config
	log     zerolog.Logger
}

func NewDBServer(port int, log zerolog.Logger, storage storage.Storage, intr *dsl.Interpreter, wal *wal.WAL) *DBServer {
	inline := exec.NewInlineExecutor(intr, wal)
	return &DBServer{
		storage: storage,
		inline:  inline,
		port:    fmt.Sprintf(":%d", port),
		log:     log,
	}
}

func (d *DBServer) Listen() error {
	ln, err := net.Listen("tcp", d.port)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go d.handleConn(conn)
	}
}

func (d *DBServer) handleConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		d.log.Debug().Msgf("got new connection line: %q", line)

		res, err := d.inline.Execute(line)
		if err != nil {
			fmt.Fprintf(conn, "%s\n", err)
			continue
		}

		d.log.Debug().Str("query", line).Msgf("got executor result data: %v", res)
		_, err = fmt.Fprintf(conn, "%v\n", res)
		if err != nil {
			break
		}
	}
}
