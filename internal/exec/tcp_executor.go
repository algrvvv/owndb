package exec

import (
	"bufio"
	"fmt"
	"net"

	"github.com/rs/zerolog"
)

type TcpExecutor struct {
	conn net.Conn
	log  zerolog.Logger
}

func NewTcpExecutor(port int, log zerolog.Logger) (Executor, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return &TcpExecutor{
		conn: conn,
		log:  log,
	}, nil
}

// Execute implements Executor.
func (t *TcpExecutor) Execute(command string) (any, error) {
	t.log.Debug().Str("command", command).Msg("got command for tcp execute")
	_, err := fmt.Fprintf(t.conn, "%s\n", command)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(t.conn)
	scanner.Scan()
	return scanner.Text(), nil
}

// Close close connection with database server
func (t *TcpExecutor) Close() error {
	return t.conn.Close()
}
