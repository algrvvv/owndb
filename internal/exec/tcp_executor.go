package exec

import (
	"bufio"
	"fmt"
	"net"
)

type TcpExecutor struct {
	conn net.Conn
}

func NewTcpExecutor(port int) (Executor, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return &TcpExecutor{
		conn: conn,
	}, nil
}

// Execute implements Executor.
func (t *TcpExecutor) Execute(command string) (any, error) {
	fmt.Println("got command for tcp execute: ", command)
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
