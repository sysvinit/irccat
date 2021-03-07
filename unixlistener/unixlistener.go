package unixlistener

import (
	"bufio"
	"github.com/irccloud/irccat/dispatcher"
	"github.com/juju/loggo"
	"github.com/spf13/viper"
	"github.com/thoj/go-ircevent"
	"net"
	"strings"
)

var log = loggo.GetLogger("UnixListener")

type UnixListener struct {
	socket net.Listener
	irc    *irc.Connection
}

func New() (*UnixListener, error) {
	var err error

	listener := UnixListener{}
	listener.socket, err = net.Listen("unix", viper.GetString("unix.listen"))
	if err != nil {
		return nil, err
	}

	return &listener, nil
}

func (l *UnixListener) Run(irccon *irc.Connection) {
	log.Infof("Listening for Unix requests on %s", viper.GetString("unix.listen"))
	l.irc = irccon
	go l.acceptConnections()
}

func (l *UnixListener) acceptConnections() {
	for {
		conn, err := l.socket.Accept()
		if err != nil {
			break
		}
		go l.handleConnection(conn)
	}
	l.socket.Close()
}

func (l *UnixListener) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.Trim(msg, "\r\n")
		if len(msg) > 0 {
			dispatcher.Send(l.irc, msg, log, conn.RemoteAddr().String())
		}
	}
	conn.Close()
}
