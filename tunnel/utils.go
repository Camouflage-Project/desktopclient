package tunnel

import (
	"github.com/function61/gokit/logex"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strings"
)

func sshClientForConn(conn net.Conn, addr string, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	sconn, chans, reqs, err := ssh.NewClientConn(conn, addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(sconn, chans, reqs), nil
}

type loggerFactory func(prefix string) *log.Logger

func mkLoggerFactory(rootLogger *log.Logger) loggerFactory {
	return func(prefix string) *log.Logger {
		return logex.Prefix(prefix, rootLogger)
	}
}

func isWebsocketAddress(address string) bool {
	return strings.HasPrefix(address, "ws://") || strings.HasPrefix(address, "wss://")
}