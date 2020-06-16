package rtmp

import (
	"net"

	"github.com/amikai/gogolive/protocol/rtmp/core"
	log "github.com/sirupsen/logrus"
)

type Server struct {
}

func NewRtmpServer() *Server {
	return &Server{}
}

func (s *Server) Serve(listener net.Listener) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("rtmp serve panic: ", r)
		}
	}()

	for {
		netConn, err := listener.Accept()
		if err != nil {
			return
		}
		rtmpConn := core.NewConn(netConn, 4*1024)
		log.Debug("New client, connect remote: ", rtmpConn.RemoteAddr().String(),
			"local:", rtmpConn.LocalAddr().String())
		go s.handleConn(rtmpConn)
	}
}

func (s *Server) handleConn(rtmpConn *core.Conn) error {
	if err := rtmpConn.SrvHandshake(); err != nil {
		rtmpConn.Close()
		log.Error("handleConn SrvHandShake err: ", err)
		return err
	}
	// TODO
	return nil
}
