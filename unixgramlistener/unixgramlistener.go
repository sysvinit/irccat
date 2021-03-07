package unixgramlistener

import (
	"bufio"
	"github.com/irccloud/irccat/dispatcher"
	"github.com/juju/loggo"
	"github.com/spf13/viper"
	"github.com/thoj/go-ircevent"
	"net"
	"strings"
)

var log = loggo.GetLogger("UnixGramListener")

type UnixGramListener struct {
	socket *net.UnixConn
	irc    *irc.Connection
}

func New() (*UnixGramListener, error) {
	var err error

	listener := UnixGramListener{}
	listener.socket, err = net.ListenUnixgram(
		"unixgram",
		&net.UnixAddr{Name: viper.GetString("unixgram.listen"), Net: "unixgram"},
	)
	if err != nil {
		return nil, err
	}

	return &listener, nil
}

func (l *UnixGramListener) Run(irccon *irc.Connection) {
	log.Infof("Listening for Unix datagram requests on %s", viper.GetString("unixgram.listen"))
	l.irc = irccon
	go l.readMessages()
}

func (l *UnixGramListener) readMessages() {
	reader := bufio.NewReader(l.socket)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.Trim(msg, "\r\n")
		if len(msg) > 0 {
			dispatcher.Send(l.irc, msg, log, "[unixgram]")
		}
	}
	l.socket.Close()
}

