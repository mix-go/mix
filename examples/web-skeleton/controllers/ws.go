package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mix-go/web-skeleton/di"
	"github.com/mix-go/xcli"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketController struct {
}

func (t *WebSocketController) Index(c *gin.Context) {
	logger := di.Logrus()
	if xcli.App().Debug {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	session := WebSocketSession{
		Conn:   conn,
		Header: c.Request.Header,
		Send:   make(chan []byte, 100),
	}
	session.Start()

	server := di.Server()
	server.RegisterOnShutdown(func() {
		session.Stop()
	})

	logger.Infof("Upgrade: %s", c.Request.UserAgent())
}

type WebSocketSession struct {
	Conn   *websocket.Conn
	Header http.Header
	Send   chan []byte
}

func (t *WebSocketSession) Start() {
	go func() {
		logger := di.Logrus()
		for {
			msgType, msg, err := t.Conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, 1001, 1006) {
					logger.Error(err)
				}
				t.Stop()
				return
			}
			if msgType != websocket.TextMessage {
				continue
			}

			handler := WebSocketHandler{
				Session: t,
			}
			handler.Index(msg)
		}
	}()
	go func() {
		logger := di.Logrus()
		for {
			msg, ok := <-t.Send
			if !ok {
				return
			}
			if err := t.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.Error(err)
				t.Stop()
				return
			}
		}
	}()
}

func (t *WebSocketSession) Stop() {
	defer func() {
		if err := recover(); err != nil {
			logger := di.Logrus()
			logger.Error(err)
		}
	}()
	close(t.Send)
	_ = t.Conn.Close()
}

type WebSocketHandler struct {
	Session *WebSocketSession
}

func (t *WebSocketHandler) Index(msg []byte) {
	t.Session.Send <- []byte("hello, world!")
}
