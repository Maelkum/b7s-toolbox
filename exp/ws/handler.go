package main

import (
	"fmt"
	"time"

	"github.com/lxzan/gws"
)

const PingInterval = 10 * time.Second

type Handler struct{}

func (c *Handler) OnOpen(socket *gws.Conn) {

	conn := socket.NetConn()
	remote := conn.RemoteAddr().String()
	local := conn.LocalAddr().String()

	log.Info().Str("local_addr", local).Str("remote_addr", remote).Msg("got new connection")

	_ = socket.SetDeadline(time.Now().Add(3 * PingInterval))

}

func (c *Handler) DeleteSession(socket *gws.Conn) {
	log.Info().Msg("delete session")
}

func (c *Handler) OnError(socket *gws.Conn, err error) {
	log.Error().Err(err).Msg("got error")
	c.DeleteSession(socket)
}

func (c *Handler) OnClose(socket *gws.Conn, code uint16, reason []byte) {
	log.Info().
		Uint16("code", code).
		Str("reason", string(reason)).
		Msg("connection closed")

	c.DeleteSession(socket)
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	log.Info().Msg("got ping")

	_ = socket.SetDeadline(time.Now().Add(3 * PingInterval))
	_ = socket.WritePong(nil)
}

func (c *Handler) OnPong(socket *gws.Conn, payload []byte) {
	log.Info().Msg("got pong")
}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	log.Info().Msg("got message")
	fmt.Printf("> %s", message.Data.String())
}
